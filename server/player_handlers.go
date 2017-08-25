package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/entities"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	RegisterOk     = "\"Player register successful.\""
	RegisterError  = "\"Player register error.\""
	PlayerNotFound = "\"Player not found.\""
	PlayerFoundOk  = "\"Player found successful.\""
)

func (env *Env) PlayerLoginPost(w http.ResponseWriter, r *http.Request) {
	session := env.getSession(w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	player := new(entities.Player)
	if success := parsePlayer(w, r, player); !success {
		return
	}

	found, err := env.PlayerDAO.FindPlayer(player)
	if found.Equal(player) && err == nil {
		session.Values[env.playerId] = found.Id
		session.Save(r, w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(PlayerFoundOk))
	} else {
		delete(session.Values, env.playerId)
		session.Save(r, w)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(PlayerNotFound))
	}
}

func (env *Env) PlayerRegisterPost(w http.ResponseWriter, r *http.Request) {
	var player entities.Player
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if success := parsePlayer(w, r, &player); !success {
		return
	}

	err := env.PlayerDAO.CreatePlayer(&player)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(RegisterOk))
	} else {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(RegisterError))
	}
}

func parsePlayer(w http.ResponseWriter, r *http.Request, player *entities.Player) (success bool) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		success = false
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		success = false
		return
	}
	if err := json.Unmarshal(body, &player); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		success = false
		return
	}
	if player.Login == "" || player.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty player"))
		success = false
		return
	}

	success = true
	return
}
