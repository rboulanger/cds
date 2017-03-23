package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ovh/cds/engine/api/worker"
	"github.com/ovh/cds/engine/log"
	"github.com/ovh/cds/sdk"
)

var (
	//VERSION is set with -ldflags "-X main.VERSION={{.cds.proj.version}}+{{.cds.version}}"
	VERSION = "snapshot"
	// WorkerID is a unique identifier for this worker
	WorkerID string
	// key is the token generated by the user owning the worker
	key         string
	name        string
	api         string
	model       int64
	hatchery    int64
	basedir     string
	bookedJobID int64
	logChan     chan sdk.Log
	// port of variable exporter HTTP server
	exportport int
	// current actionBuild is here to allow var export
	pbJob          sdk.PipelineBuildJob
	currentStep    int
	buildVariables []sdk.Variable
	// Git ssh configuration
	pkey           string
	gitsshPath     string
	startTimestamp time.Time
	nbActionsDone  int
	status         struct {
		Name      string    `json:"name"`
		Heartbeat time.Time `json:"heartbeat"`
		Status    string    `json:"status"`
		Model     int64     `json:"model"`
	}
)

var mainCmd = &cobra.Command{
	Use:   "worker",
	Short: "CDS Worker",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("cds")
		viper.AutomaticEnv()

		log.Initialize()

		log.Notice("What a good time to be alive\n")
		var err error

		name, err = os.Hostname()
		if err != nil {
			log.Notice("Cannot retrieve hostname: %s\n", err)
			return
		}

		hatchS := viper.GetString("hatchery")
		hatchery, err = strconv.ParseInt(hatchS, 10, 64)
		if err != nil {
			fmt.Printf("WARNING: Invalid hatchery ID (%s)\n", err)
		}

		api = viper.GetString("api")
		if api == "" {
			fmt.Printf("--api not provided, aborting.\n")
			return
		}

		key = viper.GetString("key")
		if key == "" {
			fmt.Printf("--key not provided, aborting.\n")
			return
		}

		givenName := viper.GetString("name")
		if givenName != "" {
			name = givenName
		}
		status.Name = name

		basedir = viper.GetString("basedir")
		if basedir == "" {
			basedir = os.TempDir()
		}

		bookedJobID = viper.GetInt64("booked_job_id")

		model = int64(viper.GetInt("model"))
		status.Model = model

		port, err := server()
		if err != nil {
			sdk.Exit("cannot bind port for worker export: %s\n", err)
		}
		exportport = port

		startTimestamp = time.Now()

		// start logger routine
		logChan = make(chan sdk.Log)
		go logger(logChan)

		go heartbeat()
		queuePolling()
	},
}

func init() {
	flags := mainCmd.Flags()

	flags.String("log-level", "notice", "Log Level : debug, info, notice, warning, critical")
	viper.BindPFlag("log_level", flags.Lookup("log-level"))

	flags.String("api", "", "URL of CDS API")
	viper.BindPFlag("api", flags.Lookup("api"))

	flags.String("key", "", "CDS KEY")
	viper.BindPFlag("key", flags.Lookup("key"))

	flags.Bool("single-use", false, "Exit after executing an action")
	viper.BindPFlag("single_use", flags.Lookup("single-use"))

	flags.String("name", "", "Name of worker")
	viper.BindPFlag("name", flags.Lookup("name"))

	flags.Int("model", 0, "Model of worker")
	viper.BindPFlag("model", flags.Lookup("model"))

	flags.Int("hatchery", 0, "Hatchery spawing worker")
	viper.BindPFlag("hatchery", flags.Lookup("hatchery"))

	flags.String("basedir", "", "Worker working directory")
	viper.BindPFlag("basedir", flags.Lookup("basedir"))

	flags.Int("ttl", 30, "Worker time to live (minutes)")
	viper.BindPFlag("ttl", flags.Lookup("ttl"))

	flags.Int("heartbeat", 10, "Worker heartbeat frequency")
	viper.BindPFlag("heartbeat", flags.Lookup("heartbeat"))

	flags.Int64("booked-job-id", 0, "Booked job id")
	viper.BindPFlag("booked_job_id", flags.Lookup("booked-job-id"))

	mainCmd.AddCommand(cmdExport)
	mainCmd.AddCommand(cmdUpload)
	mainCmd.AddCommand(versionCmd)
}

func main() {
	sdk.SetAgent(sdk.WorkerAgent)

	mainCmd.Execute()
}

// Will be removed when websocket conn is implemented
// for now, poll the /queue
func queuePolling() {
	firstViewQueue := true
	for {
		if WorkerID == "" {
			log.Notice("[WORKER] Registering on CDS engine, \n")
			if err := register(api, name, key); err != nil {
				log.Notice("Cannot register: %s\n", err)
				time.Sleep(10 * time.Second)
				continue
			}
		}

		//We we've done nothing until ttl is over, let's exit
		if nbActionsDone == 0 && startTimestamp.Add(time.Duration(viper.GetInt("ttl"))*time.Minute).Before(time.Now()) {
			log.Notice("Time to exit.")
			unregister()
			os.Exit(0)
		}

		checkQueue(bookedJobID)
		if firstViewQueue {
			// if worker did not found booked job ID is first iteration
			// reset booked job to take another action
			bookedJobID = 0
		}

		time.Sleep(5 * time.Second)
	}
}

func checkQueue(bookedJobID int64) {
	defer sdk.SetWorkerStatus(sdk.StatusWaiting)

	queue, err := sdk.GetBuildQueue()
	if err != nil {
		log.Notice("checkQueue> Cannot get build queue: %s\n", err)
		time.Sleep(5 * time.Second)
		WorkerID = ""
		return
	}

	log.Notice("checkQueue> %d Actions in queue", len(queue))

	//Set the status to checking to avoid beeing killed while checking queue, actions and requirements
	sdk.SetWorkerStatus(sdk.StatusChecking)

	for i := range queue {
		if bookedJobID != 0 && queue[i].ID != bookedJobID {
			continue
		}

		requirementsOK := true
		// Check requirement
		log.Notice("checkQueue> Checking requirements for action [%d] %s", queue[i].ID, queue[i].Job.Action.Name)
		for _, r := range queue[i].Job.Action.Requirements {
			ok, err := checkRequirement(r)
			if err != nil {
				postCheckRequirementError(&r, err)
				requirementsOK = false
				continue
			}
			if !ok {
				requirementsOK = false
				continue
			}
		}

		if requirementsOK {
			t := ""
			if queue[i].ID != bookedJobID {
				t = ", this was my booked job"
			}
			log.Notice("checkQueue> Taking job %d%s", queue[i].ID, t)
			takeAction(queue[i], queue[i].ID == bookedJobID)
		}
	}

	log.Notice("checkQueue> Nothing to do...")
}

func postCheckRequirementError(r *sdk.Requirement, err error) {
	s := fmt.Sprintf("Error checking requirement Name=%s Type=%s Value=%s :%s", r.Name, r.Type, r.Value, err)
	btes := []byte(s)
	sdk.Request("POST", "/queue/requirements/errors", btes)
}

func takeAction(b sdk.PipelineBuildJob, isBooked bool) {
	in := worker.TakeForm{Time: time.Now()}
	if isBooked {
		in.BookedJobID = b.ID
	}

	bodyTake, errm := json.Marshal(in)
	if errm != nil {
		log.Notice("takeAction: Cannot marshal body: %s\n", errm)
	}

	nbActionsDone++
	gitsshPath = ""
	pkey = ""
	path := fmt.Sprintf("/queue/%d/take", b.ID)
	data, code, errr := sdk.Request("POST", path, bodyTake)
	if errr != nil {
		log.Notice("takeAction> Cannot take action %d : %s\n", b.Job.PipelineActionID, errr)
		return
	}
	if code != http.StatusOK {
		return
	}

	pbji := worker.PipelineBuildJobInfo{}
	if err := json.Unmarshal([]byte(data), &pbji); err != nil {
		log.Notice("takeAction> Cannot unmarshal action: %s\n", err)
		return
	}

	pbJob = pbji.PipelineBuildJob
	// Reset build variables
	buildVariables = nil
	start := time.Now()
	res := run(&pbji)
	res.RemoteTime = time.Now()
	res.Duration = sdk.Round(time.Since(start), time.Second).String()

	// Give time to buffered logs to be sent
	time.Sleep(3 * time.Second)

	path = fmt.Sprintf("/queue/%d/result", b.ID)
	body, errm := json.MarshalIndent(res, " ", " ")
	if errm != nil {
		log.Notice("takeAction>Cannot marshal result: %s\n", errm)
		return
	}

	code = 300
	var isThereAnyHopeLeft = 10
	for code >= 300 {
		var errre error
		_, code, errre = sdk.Request("POST", path, body)
		if code == http.StatusNotFound {
			unregister() // well...
			log.Notice("takeAction> Cannot send build result: PipelineBuildJob does not exists anymore\n")
			break
		}
		if errre == nil && code < 300 {
			fmt.Printf("BuildResult sent.\n")
			break
		}

		if errre != nil {
			log.Notice("takeAction> Cannot send build result: %s\n", errre)
		} else {
			log.Notice("takeAction> Cannot send build result: HTTP %d\n", code)
		}

		time.Sleep(5 * time.Second)
		isThereAnyHopeLeft--
		if isThereAnyHopeLeft < 0 {
			log.Notice("takeAction> Could not send built result 50 times, giving up\n")
			break
		}
	}

	if viper.GetBool("single_use") {
		// Give time to logs to be flushed
		time.Sleep(2 * time.Second)
		// Unregister from engine
		if err := unregister(); err != nil {
			log.Warning("takeAction> could not unregister: %s\n", err)
		}
		// then exit
		log.Notice("takeAction> --single_use is on, exiting\n")
		os.Exit(0)
	}

}

func heartbeat() {
	for {
		time.Sleep(time.Duration(viper.GetInt("heartbeat")) * time.Second)
		if WorkerID == "" {
			log.Notice("[WORKER] Disconnected from CDS engine, trying to register...\n")
			if err := register(api, name, key); err != nil {
				log.Notice("Cannot register: %s\n", err)
				continue
			}
		}

		_, code, err := sdk.Request("POST", "/worker/refresh", nil)
		if err != nil || code >= 300 {
			log.Notice("heartbeat> cannot refresh beat: %d %s\n", code, err)
			WorkerID = ""
		}
	}
}

func unregister() error {
	uri := "/worker/unregister"
	_, code, err := sdk.Request("POST", uri, nil)
	if err != nil {
		return err
	}
	if code > 300 {
		return fmt.Errorf("HTTP %d", code)
	}
	return nil
}
