package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/management"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/rentgen94/QuestGoMail/entities"
)

const (
	alreadyPlaying = "You are already playing"
	notPlaying = "You are not playing yet"

	bagCapacity = 1000
	inputPlayerBuffSize = 10
	outputPlayerBuffSize = 10
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
	env.gameCommandGet(w, r, management.GetRoomCode)
}

func (env *Env) GameSlotsGet(w http.ResponseWriter, r *http.Request) {
	env.gameCommandGet(w, r, management.GetSlotsCode)
}

func (env *Env) GameBagGet(w http.ResponseWriter, r *http.Request) {
	env.gameCommandGet(w, r, management.GetBagCode)
}

func (env *Env) GameDoorsGet(w http.ResponseWriter, r *http.Request) {
	env.gameCommandGet(w, r, management.GetDoorsCode)
}

func (env *Env) GameItemsGet(w http.ResponseWriter, r *http.Request) {
	env.gameCommandGet(w, r, management.GetItemsCode)
}

func (env *Env) GameInteractivesGet(w http.ResponseWriter, r *http.Request) {
	env.gameCommandGet(w, r, management.GetIteractivesCode)
}

func (env *Env) GameStartPost(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var labIdStr = vars["labyrinth_id"]
	var labId, err = strconv.Atoi(labIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	var session = env.getSession(w, r)
	if !env.authorized(session) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if env.playing(session) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(alreadyPlaying))
	}

	var labyrinth, labErr = env.LabyrinthDao.GetById(labId)
	if labErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(labErr.Error()))
	}

	var player = entities.NewPlayer(labyrinth, bagCapacity)
	var manager, managerErr = management.NewPlayerManager(player, inputPlayerBuffSize, outputPlayerBuffSize)

	if managerErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(managerErr.Error()))
	}

	var gameId = env.Pool.AddManager(manager)
	session.Values[env.gameId] = gameId
	session.Save(r, w)
}

func (env *Env) GameQuitPost(w http.ResponseWriter, r *http.Request) {
	var session = env.getSession(w, r)
	if !env.authorized(session) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if env.playing(session) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(notPlaying))
	}

	var gameId = session.Values[env.gameId].(int)
	env.Pool.DeleteManager(gameId)
	delete(session.Values, env.gameId)
	session.Save(r, w)
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

func (env *Env) gameCommandGet(w http.ResponseWriter, r *http.Request, commandType int) {
	command := management.NewCommand(commandType, 0, nil, nil)
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
