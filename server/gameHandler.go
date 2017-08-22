package server

import (
	"encoding/json"
	"github.com/rentgen94/QuestGoMail/management"
	"net/http"
	"github.com/rentgen94/QuestGoMail/entities"
)

func (env *Env) GameLookAroundGet (w http.ResponseWriter, r *http.Request) {
	env.GameComandGet (w, r, management.GetRoomCode)
}

func (env *Env)  GameSlotsGet (w http.ResponseWriter, r *http.Request) {
	env.GameComandGet (w, r, management.GetSlotsCode)
}


func (env *Env)  GameBagGet (w http.ResponseWriter, r *http.Request) {
	env.GameComandGet (w, r, management.GetBagCode)
}


func (env *Env)  GameDoorsGet (w http.ResponseWriter, r *http.Request) {
	env.GameComandGet (w, r, management.GetDoorsCode)
}


func (env *Env)  GameItemsGet (w http.ResponseWriter, r *http.Request) {
	env.GameComandGet (w, r, management.GetItemsCode)
}


func (env *Env)  GameIteractivesGet (w http.ResponseWriter, r *http.Request) {
	env.GameComandGet (w, r, management.GetIteractivesCode)
}


func (env *Env)  GameComandGet(w http.ResponseWriter, r *http.Request, comandType int) {
	session := env.getSession(w, r)

	if session.Values[env.playerId] == nil {
		// Игрок не авторизован
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
	}
	command := management.NewCommand (comandType, "", nil, nil)

	//Todo Привязать к основному PoolManager
	managerPool := management.NewManagerPool(10, 10)
	//Todo получать game_id
	game_id, _ := session.Values[env.playerId].(int)
	room := entities.NewRoom(0, "ny first room", "our demons hide in the dark")
	player := entities.NewPlayer(room, 100)
	var manager, _ = management.NewPlayerManager(player, 10, 10)
	managerPool.AddManager(manager)
	managerPool.SendCommand(management.AddressedCommand{game_id, command})
	response := managerPool.GetResponseSync(game_id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(res)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
