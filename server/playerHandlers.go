package server

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/rentgen94/QuestGoMail/server/database"
)

func PlayerLoginPost(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player := new(entities.Player)
	var msg string
	parsePlayer(w, r, player)

	if id := session.Values[authToken]; id != nil {
		if founded, _ := database.FindPlayerById(id.(int));
			founded != nil && founded.Login == player.Login && founded.Password == player.Password  {
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
			// Set some session values.
			session.Values[authToken] = player.Id
			// Save it before we write to the response/return from the handler.
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
	if msg == database.RegisterOk {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	} else if msg == database.AlreadyRegistered {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		return
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
