package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/rentgen94/QuestGoMail/server/database"
	"io"
	"io/ioutil"
	"net/http"
)

func (env *Env) PlayerLoginPost(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	player := new(entities.Player)
	var msg string
	parsePlayer(w, r, player)

	if id := session.Values[authToken]; id != nil {
		if founded, _ := env.PlayerDAO.FindPlayerById(id.(int)); founded != nil && founded.Login == player.Login && founded.Password == player.Password {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			session.Values[authToken] = nil
			session.Save(r, w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		player, msg = env.PlayerDAO.FindPlayer(player)

		if msg == database.PlayerFoundOk {
			session.Values[authToken] = player.Id
			session.Save(r, w)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (env *Env) PlayerRegisterPost(w http.ResponseWriter, r *http.Request) {
	var player entities.Player
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	parsePlayer(w, r, &player)

	msg := env.PlayerDAO.CreatePlayer(&player)
	switch msg {
	case database.RegisterOk:
		w.WriteHeader(http.StatusCreated)
	case database.AlreadyRegistered:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parsePlayer(w http.ResponseWriter, r *http.Request, player *entities.Player) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &player); err != nil {
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
}
