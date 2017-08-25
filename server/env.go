package server

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/rentgen94/QuestGoMail/database/dao"
	"github.com/rentgen94/QuestGoMail/management"
	"net/http"
)

const (
	maxReadBytes = 1048576
)

type Env struct {
	PlayerDAO  dao.PlayerDAO
	Store      *sessions.CookieStore
	Pool       *management.ManagerPool
	playerId   string
	cookieName string
	gameId     string
	curGame    int
}

func NewEnv(db *sql.DB, poolWorkerNum int, commandBuffSize int, respBuffSize int) *Env {
	return &Env{
		PlayerDAO:  dao.NewDBPlayerDAO(db),
		Store:      sessions.NewCookieStore([]byte("server-cookie-store")),
		Pool:       management.NewManagerPool(poolWorkerNum, commandBuffSize, respBuffSize),
		playerId:   "player_id",
		gameId:     "game_id",
		cookieName: "quest_go_mail",
		curGame:    1,
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
