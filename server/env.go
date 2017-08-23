package server

import (
	"github.com/gorilla/sessions"
	"github.com/rentgen94/QuestGoMail/server/database"
	"net/http"
)

type Env struct {
	PlayerDAO  database.PlayerDAO
	Store      *sessions.CookieStore
	authToken  string
	cookieName string
}

func NewEnv() Env {
	return Env{
		PlayerDAO:  database.NewDBPlayerDAO(database.Init()),
		Store:      sessions.NewCookieStore([]byte("server-cookie-store")),
		authToken:  "auth_token",
		cookieName: "quest_go_mail",
	}
}

func (env *Env) getSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	// Get a session. Get() always returns a session, even if empty.
	session, err := env.Store.Get(r, env.cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return session
}
