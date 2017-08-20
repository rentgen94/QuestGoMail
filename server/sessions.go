package server

import (
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	sessionName = "Session-for-new-user"
	authToken   = "auth_token"
)

var store = sessions.NewCookieStore([]byte("server-cookie-store"))

func getSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return session
}
