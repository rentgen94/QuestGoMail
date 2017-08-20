package server

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"github.com/rentgen94/QuestGoMail/entities"
	"github.com/rentgen94/QuestGoMail/server/database"
	"fmt"
)

func PlayerLoginPost(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("this is party id", session.Values["party_id"])

	if id := session.Values["party_id"]; id != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		return
	} else {
		var player entities.Player
		parsePlayer(w, r, &player)

		_, msg := database.FindPlayer(&player)

		if msg == database.PlayerFoundOk {

			fmt.Println("this is player id", player.Id)
			// Set some session values.
			session.Values["party_id"] = player.Id
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
