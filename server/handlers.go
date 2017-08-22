package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/management"
	"io"
	"io/ioutil"
	"net/http"
)

type Default struct {
}

func (env *Env) GameCommandPost(w http.ResponseWriter, r *http.Request) {
	session := env.getSession(w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if session.Values[env.authToken] == nil {
		// Игрок не авторизован
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var command management.Command

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &command); err != nil {
		w.WriteHeader(http.StatusBadRequest) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	//Todo Привязать к основному PoolManager
	managerPool := management.NewManagerPool(10, 10)
	//Todo получать game_id
	game_id, _ := session.Values[env.authToken].(int)
	managerPool.SendCommand(management.AddressedCommand{game_id, command})
	response := managerPool.GetResponseSync(game_id)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (env *Env) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}