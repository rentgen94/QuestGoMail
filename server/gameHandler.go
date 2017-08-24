package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/management"
	"net/http"
)

func (env *Env) GameLookAroundGet(w http.ResponseWriter, r *http.Request) {
	env.GameComandGet(w, r, management.GetRoomCode)
}

func (env *Env) GameSlotsGet(w http.ResponseWriter, r *http.Request) {
	env.GameComandGet(w, r, management.GetSlotsCode)
}

func (env *Env) GameBagGet(w http.ResponseWriter, r *http.Request) {
	env.GameComandGet(w, r, management.GetBagCode)
}

func (env *Env) GameDoorsGet(w http.ResponseWriter, r *http.Request) {
	env.GameComandGet(w, r, management.GetDoorsCode)
}

func (env *Env) GameItemsGet(w http.ResponseWriter, r *http.Request) {
	env.GameComandGet(w, r, management.GetItemsCode)
}

func (env *Env) GameIteractivesGet(w http.ResponseWriter, r *http.Request) {
	env.GameComandGet(w, r, management.GetIteractivesCode)
}

func (env *Env) GameComandGet(w http.ResponseWriter, r *http.Request, comandType int) {
	session := env.getSession(w, r)

	if session.Values[env.playerId] == nil {
		// Игрок не авторизован
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
	}

	if session.Values[env.gameId] == nil {
		// Игрок не в игре
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	command := management.NewCommand (comandType, 0, nil, nil)

	env.Pool.SendCommand(management.AddressedCommand{session.Values[env.gameId].(int), command})
	response, err := env.Pool.GetResponseSync(session.Values[env.gameId].(int))
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeInternalError(w)
		return
	}
}
