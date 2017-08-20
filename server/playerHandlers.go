package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/rentgen94/QuestGoMail/server/database"
	"io"
	"io/ioutil"
	"net/http"
)

func PlayerLoginPost(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)

	player := new(entities.Player)
	var msg string
	parsePlayer(w, r, player)

	if id := session.Values[authToken]; id != nil {
		if founded, _ := database.FindPlayerById(id.(int)); founded != nil && founded.Login == player.Login && founded.Password == player.Password {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			session.Values[authToken] = nil
			session.Save(r, w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		player, msg = database.FindPlayer(player)

		if msg == database.PlayerFoundOk {
			session.Values[authToken] = player.Id
			session.Save(r, w)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func PlayerRegisterPost(w http.ResponseWriter, r *http.Request) {
	var player entities.Player
	parsePlayer(w, r, &player)

	msg := database.CreateUser(&player)
	switch msg {
	case database.RegisterOk:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	case database.AlreadyRegistered:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
}

