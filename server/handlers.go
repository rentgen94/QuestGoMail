package server

import (
	"net/http"
	"github.com/rentgen94/QuestGoMail/management"
	"io/ioutil"
	"io"
	"encoding/json"
)

type Default struct {

}

func GameCommandPost(w http.ResponseWriter, r *http.Request) {
	// Получаю
	session, err := Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["party_id"] != nil {
		var command management.Command

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &command); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}

		//Todo Привязать к основному PoolManager
		//managerPool := management.NewManagerPool();
		//party_id, _ := strconv.Atoi(string(session.Values["party_id"]))
		//managerPool.SendCommand(management.AddressedCommand{party_id, command})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		// Игрок не авторизован
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
	}
}