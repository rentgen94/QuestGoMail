package server

import (
	"github.com/rentgen94/QuestGoMail/server/database"
	"github.com/gorilla/sessions"
	"net/http"
	"github.com/rentgen94/QuestGoMail/management"
)

const (
	maxReadBytes = 1048576
)

type Env struct {
	PlayerDAO  database.PlayerDAO
	Store      *sessions.CookieStore
	Pool       *management.ManagerPool
	playerId   string
	cookieName string
}

func NewEnv() Env {
	return Env{
		PlayerDAO:  database.NewDBPlayerDAO(database.Init()),
		Store:      sessions.NewCookieStore([]byte("server-cookie-store")),
		playerId:   "player_id",
		cookieName: "quest_go_mail",
		Pool: management.NewManagerPool(10, 10),
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

func writeInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
