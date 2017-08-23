package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/entities"
	"io"
	"io/ioutil"
	"net/http"
	"errors"
	"github.com/rentgen94/QuestGoMail/management"
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
	if err := parsePlayer(w, r, player); err != nil {
		writeInternalError(w)
		return
	}

	founded, err := env.PlayerDAO.FindPlayer(player)
	if equal(founded, player) && err == nil {
		session.Values[env.playerId] = founded.Id
		session.Save(r, w)
		room := entities.NewRoom(0, "my first room", "our demons hide in the dark")
		*player.Room() = *room
		var manager, _ = management.NewPlayerManager(player)
		env.Pool.AddManager(manager)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(PlayerFoundOk))
		return
	} else {
		session.Values[env.playerId] = nil
		session.Save(r, w)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(PlayerNotFound))
		return
	}
}

func (env *Env) PlayerRegisterPost(w http.ResponseWriter, r *http.Request) {
	var player entities.Player
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := parsePlayer(w, r, &player); err != nil {
		writeInternalError(w)
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

func parsePlayer(w http.ResponseWriter, r *http.Request, player *entities.Player) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, &player); err != nil {
		return err
	}
	if player.Login == "" || player.Password == "" {
		return errors.New("Empty player")
	}

	return nil
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
