package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/management"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Default struct {
}

func (env *Env) GameCommandPost(w http.ResponseWriter, r *http.Request) {
	session := env.getSession(w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if session.IsNew {
		// Игрок не авторизован
		w.WriteHeader(http.StatusForbidden)
		return
	}
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
		w.WriteHeader(http.StatusBadRequest) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			writeInternalError(w)
			return
		}
		return
	}

	if session.Values[env.gameId] == nil {
		// Игрок не в игре
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	env.Pool.SendCommand(management.AddressedCommand{session.Values[env.gameId].(int), command})
	response, err := env.Pool.GetResponseSync(session.Values[env.gameId].(int), time.Minute)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeInternalError(w)
		return
	}
}

func (env *Env) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}
