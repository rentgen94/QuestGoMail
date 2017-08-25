package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/management"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/gorilla/sessions"
)

func (env *Env) GameListLabyrinthsGet(w http.ResponseWriter, r *http.Request) {
	var labs, err = env.LabyrinthDao.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(labs); err != nil {
		writeInternalError(w)
		return
	}
}

func (env *Env) GameLookAroundGet(w http.ResponseWriter, r *http.Request) {
	env.GameCommandGet(w, r, management.GetRoomCode)
}

func (env *Env) GameSlotsGet(w http.ResponseWriter, r *http.Request) {
	env.GameCommandGet(w, r, management.GetSlotsCode)
}

func (env *Env) GameBagGet(w http.ResponseWriter, r *http.Request) {
	env.GameCommandGet(w, r, management.GetBagCode)
}

func (env *Env) GameDoorsGet(w http.ResponseWriter, r *http.Request) {
	env.GameCommandGet(w, r, management.GetDoorsCode)
}

func (env *Env) GameItemsGet(w http.ResponseWriter, r *http.Request) {
	env.GameCommandGet(w, r, management.GetItemsCode)
}

func (env *Env) GameInteractivesGet(w http.ResponseWriter, r *http.Request) {
	env.GameCommandGet(w, r, management.GetIteractivesCode)
}

func (env *Env) GameCommandGet(w http.ResponseWriter, r *http.Request, commandType int) {
	command := management.NewCommand(commandType, 0, nil, nil)
	env.handleGameCommand(w, r, command)
}

func (env *Env) GameCommandPost(w http.ResponseWriter, r *http.Request) {
	var command management.Command

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
	if err != nil {
		writeInternalError(w)
		return
	}
	if err := r.Body.Close(); err != nil {
		writeInternalError(w)
		return
	}
	if err := json.Unmarshal(body, &command); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			writeInternalError(w)
			return
		}
		return
	}

	env.handleGameCommand(w, r, command)
}

func (env *Env) handleGameCommand(w http.ResponseWriter, r *http.Request, command management.Command) {
	session := env.getSession(w, r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if !env.authorized(session) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !env.playing(session) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	env.Pool.SendCommand(
		management.AddressedCommand{
			session.Values[env.gameId].(int),
			command,
		},
	)
	response, err := env.Pool.GetResponseSync(session.Values[env.gameId].(int), 1*time.Second)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeInternalError(w)
		return
	}
}

func (env *Env) authorized(session *sessions.Session) bool {
	return !session.IsNew
}

func (env *Env) playing(session *sessions.Session) bool {
	return session.Values[env.gameId] != nil
}
