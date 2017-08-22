package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/entities"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	RegisterOk        = "\"Player register successful.\""
	RegisterError     = "\"Player register error.\""
	PlayerNotFound    = "\"Player not found.\""
	PlayerFoundOk     = "\"Player found successful.\""
)

func (env *Env) PlayerLoginPost(w http.ResponseWriter, r *http.Request) {
	session:= env.getSession(w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	player := new(entities.Player)
	parsePlayer(w, r, player)

	founded, err := env.PlayerDAO.FindPlayer(player)
	if equal(founded, player) && err == nil {
		session.Values[env.authToken] = founded.Id
		session.Save(r, w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(PlayerFoundOk))
		return
	} else {
		session.Values[env.authToken] = nil
		session.Save(r, w)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(PlayerNotFound))
		return
	}
}

func (env *Env) PlayerRegisterPost(w http.ResponseWriter, r *http.Request) {
	var player entities.Player
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	parsePlayer(w, r, &player)

	err := env.PlayerDAO.CreatePlayer(&player)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(RegisterOk))
	} else {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(RegisterError))
	}
}

func parsePlayer(w http.ResponseWriter, r *http.Request, player *entities.Player) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &player); err != nil {
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func equal(this, other *entities.Player) (bool) {
	if this == nil || other == nil {
		return false
	} else if this.Login != other.Login {
		return false
	} else if this.Password != other.Password {
		return false
	}

	return true
}
