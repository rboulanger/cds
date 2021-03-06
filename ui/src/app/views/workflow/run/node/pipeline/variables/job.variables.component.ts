import {Component, Input, ViewChild} from '@angular/core';
import {SemanticModalComponent} from 'ng-semantic/ng-semantic';
import {Variable} from '../../../../../../model/variable.model';
import {Parameter} from '../../../../../../model/parameter.model';

@Component({
    selector: 'app-workflow-run-job-variable',
    templateUrl: './job.variable.html',
    styleUrls: ['./job.variable.scss']
})
export class WorkflowRunJobVariableComponent {

    @ViewChild('jobVariablesModal')
    jobVariablesModal: SemanticModalComponent;

    @Input('variables')
    set variables(data: Array<Variable>) {
        this.init();
        if (data) {
            data.forEach(p => {
                if (p.name.indexOf('cds.proj.', 0) === 0) {
                    this.varProject.push(p);
                } else if (p.name.indexOf('cds.app.', 0) === 0) {
                    this.varApplication.push(p);
                } else if (p.name.indexOf('cds.pip.', 0) === 0) {
                    this.varPipeline.push(p);
                } else if (p.name.indexOf('cds.env.', 0) === 0) {
                    this.varEnvironment.push(p);
                } else if (p.name.indexOf('cds.parent.', 0) === 0) {
                    this.varParent.push(p);
                } else if (p.name.indexOf('cds.build.', 0) === 0) {
                    this.varBuild.push(p);
                } else if (p.name.indexOf('git.', 0) === 0) {
                    this.varGit.push(p);
                } else if (p.name.indexOf('workflow.', 0) === 0) {
                    this.varWorkflow.push(p);
                } else {
                    this.varCDS.push(p);
                }
            });
        }
    }

    varGit: Array<Parameter>;
    varCDS: Array<Parameter>;
    varBuild: Array<Parameter>;
    varEnvironment: Array<Parameter>;
    varApplication: Array<Parameter>;
    varPipeline: Array<Parameter>;
    varProject: Array<Parameter>;
    varParent: Array<Parameter>;
    varWorkflow: Array<Parameter>;


    constructor() {
    }

    init(): void {
        this.varGit = new Array<Parameter>();
        this.varCDS = new Array<Parameter>();
        this.varBuild = new Array<Parameter>();
        this.varEnvironment = new Array<Parameter>();
        this.varApplication = new Array<Parameter>();
        this.varPipeline = new Array<Parameter>();
        this.varProject = new Array<Parameter>();
        this.varParent = new Array<Parameter>();
        this.varWorkflow = new Array<Parameter>();
    }

    show(data: {}): void {
        if (this.jobVariablesModal) {
            this.jobVariablesModal.show(data);
        }
    }
}
