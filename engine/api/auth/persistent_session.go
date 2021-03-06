package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/go-gorp/gorp"

	"github.com/ovh/cds/engine/api/sessionstore"
	"github.com/ovh/cds/engine/api/user"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/log"
)

//NewPersistentSession create a new session with token stored as user_key in database
func NewPersistentSession(db gorp.SqlExecutor, d Driver, u *sdk.User) (sessionstore.SessionKey, error) {
	u, errLoad := user.LoadUserAndAuth(db, u.Username)
	if errLoad != nil {
		return "", errLoad
	}

	t, errSession := user.NewPersistentSession(db, u)
	if errSession != nil {
		return "", errSession
	}

	session, errStore := d.Store().New(t)
	if errStore != nil {
		return "", errStore
	}
	log.Info("NewPersistentSession> New Session for %s", u.Username)
	d.Store().Set(session, "username", u.Username)

	return session, nil
}

//CheckPersistentSession check persistent session token from CLI
func checkPersistentSession(ctx context.Context, db gorp.SqlExecutor, store sessionstore.Store, headers http.Header) (context.Context, bool) {
	if headers.Get(sdk.RequestedWithHeader) == sdk.RequestedWithValue {
		return getUserPersistentSession(ctx, db, store, headers)
	}
	return ctx, false
}

func getUserPersistentSession(ctx context.Context, db gorp.SqlExecutor, store sessionstore.Store, headers http.Header) (context.Context, bool) {
	h := headers.Get(sdk.SessionTokenHeader)
	if h == "" {
		return ctx, false
	}

	key := sessionstore.SessionKey(h)
	ok, _ := store.Exists(key)
	var err error
	var u *sdk.User

	if !ok {
		//Reload the persistent session from the database
		token, err := user.LoadPersistentSessionToken(db, key)
		if err != nil {
			log.Warning("getUserPersistentSession> Unable to load user by token %s", key)
			return ctx, false
		}
		u, err = user.LoadUserWithoutAuthByID(db, token.UserID)
		store.New(key)
		store.Set(key, "username", u.Username)

	} else {
		//The session is in the session store
		var usr string
		store.Get(sessionstore.SessionKey(h), "username", &usr)
		u, err = user.LoadUserWithoutAuth(db, usr)
	}

	//Check previous errors
	if err != nil {
		log.Warning("getUserPersistentSession> Unable to load user")
		return ctx, false
	}

	//Set user in ctx
	ctx = context.WithValue(ctx, ContextUser, u)

	//Launch update of the persistent session token in background
	defer func(key sessionstore.SessionKey) {
		token, err := user.LoadPersistentSessionToken(db, key)
		if err != nil {
			log.Warning("getUserPersistentSession> Unable to load user by token %s: %v", key, err)
			return
		}
		token.LastConnectionDate = time.Now()
		if err := user.UpdatePersistentSessionToken(db, *token); err != nil {
			log.Error("getUserPersistentSession> Unable to update token")
		}
	}(key)

	return ctx, true
}
