package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionAuthState = "authenticated"
	SessionUserID    = "currentUserId"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func requestSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "cookie-name")
	return session
}

func CurrentUserId(r *http.Request) *int64 {
	session := requestSession(r)

	userID, ok := session.Values[SessionUserID].(int64)
	if !ok {
		return nil
	}

	return &userID
}

func IsUserLoggedIn(r *http.Request) bool {
	session := requestSession(r)

	auth, ok := session.Values[SessionAuthState].(bool)

	return ok && auth
}

func SetLoggedInUserID(w http.ResponseWriter, r *http.Request, userID int64) {
	session := requestSession(r)
	session.Values[SessionAuthState] = true
	session.Values[SessionUserID] = userID
	session.Save(r, w)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session := requestSession(r)
	session.Values[SessionAuthState] = false
	session.Values[SessionUserID] = nil

	session.Save(r, w)
}
