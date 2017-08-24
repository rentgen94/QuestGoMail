package server

import (
	"github.com/gorilla/sessions"
	"github.com/rentgen94/QuestGoMail/management"
	"github.com/rentgen94/QuestGoMail/server/database"
	"net/http"
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
	gameId     string
	curGame    int
}

func NewEnv() Env {
	return Env{
		PlayerDAO:  database.NewDBPlayerDAO(database.Init()),
		Store:      sessions.NewCookieStore([]byte("server-cookie-store")),
		playerId:   "player_id",
		gameId:     "game_id",
		cookieName: "quest_go_mail",
		curGame:    1,
		Pool:       management.NewManagerPool(1, 10, 10),
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


func (env *Env) NewGame() int {
	env.curGame += 1
	return env.curGame
}