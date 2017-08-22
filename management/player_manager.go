package management

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"strconv"
)

const (
	playerInTheVoid             = "Player is not in a room"
	badCode                     = "Bad code"
	doorNotFoundTemplate        = "Door \"%s\" not found"
	interactiveNotFoundTemplate = "Interactive \"%s\" not found"
	itemCodeNotSupplied         = "Item code not supplied"
	cant_create_json            = "can't get data"

	managerNotStartedCode = iota
	managerWorkCode
	managerFinishedCode
)

type PlayerManager struct {
	stateCode int
	player    *entities.Player
	inChan    chan Command
	outChan   chan Response
	stopChan  chan interface{}
}

func NewPlayerManager(player *entities.Player, commandBuffSize int, responseBuffSize int) (*PlayerManager, error) {
	if player.Room() == nil {
		return nil, errors.New(playerInTheVoid)
	}

	return &PlayerManager{
		stateCode: managerNotStartedCode,
		player:    player,
		inChan:    make(chan Command, commandBuffSize),
		outChan:   make(chan Response, responseBuffSize),
		stopChan:  make(chan interface{}, 1),
	}, nil
}

func (manager *PlayerManager) CommandChan() chan Command {
	return manager.inChan
}

func (manager *PlayerManager) RespChan() chan Response {
	return manager.outChan
}

func (manager *PlayerManager) Run() {
	manager.stateCode = managerWorkCode

	for {
		select {
		case command := <-manager.inChan:
			var resp = manager.getCommandResponse(command)
			manager.outChan <- resp
			if resp.IsFinish() {
				break
			}
		case <-manager.stopChan:
			break
		default:
		}
	}

	manager.stateCode = managerFinishedCode
}

func (manager *PlayerManager) Stop() {
	manager.stateCode = managerFinishedCode
	manager.stopChan <- 1
}

func (manager *PlayerManager) Finished() bool {
	return manager.stateCode == managerFinishedCode
}

func (manager *PlayerManager) getCommandResponse(command Command) Response {
	var methodMap = map[int]func(*Response, *PlayerManager, Command){
		GetRoomCode:        handleRoomCode,
		GetSlotsCode:       handleSlotsCode,
		GetDoorsCode:       handleDoorsCode,
		GetItemsCode:       handleItemsCode,
		GetBagCode:         handleBagCode,
		GetIteractivesCode: handleIteractivesCode,
		enterCode:          handleEnterCode,
		interactCode:       handleInteractCode,
		takeCode:           handleTakeCode,
		putCode:            handlePutCode,
	}

	var f, ok = methodMap[command.typeCode]

	var resp = NewResponse()
	if !ok {
		resp.ErrMsg = badCode
		return resp
	}
	f(&resp, manager, command)

	return resp
}

type roomResponse struct {
	name        string `json:"name"`
	description string `json:"description"`
}

func handleRoomCode(resp *Response, manager *PlayerManager, command Command) {
	a := &roomResponse{
		name:        manager.player.Room().Name(),
		description: manager.player.Room().Description(),
	}
	res, err := json.Marshal(a)
	if err != nil {
		resp.ErrMsg = cant_create_json
		return
	}
	resp.Data = res
}

type slotsResponse struct {
	id       int    `json:"id"`
	name     string `json:"name"`
	capacity int    `json:"capacity"`
	contains int    `json:"contains"`
}

func handleSlotsCode(resp *Response, manager *PlayerManager, command Command) {
	a := []slotsResponse{}
	for _, elem := range manager.player.Room().Slots() {
		slt := &slotsResponse{
			name:     elem.Name(),
			capacity: elem.Capacity(),
			contains: elem.Contains(),
			id:       elem.Id(),
		}
		a = append(a, *slt)
	}
	res, err := json.Marshal(a)
	if err != nil {
		resp.ErrMsg = cant_create_json
		return
	}
	resp.Data = res
}

type doorsResponse struct {
	id   int    `json:"id"`
	name string `json:"name"`
}

func handleDoorsCode(resp *Response, manager *PlayerManager, command Command) {
	a := []doorsResponse{}
	for _, elem := range manager.player.Room().Doors() {
		slt := &doorsResponse{
			name: elem.Name(),
			id:   elem.Id(),
		}
		a = append(a, *slt)
	}
	res, err := json.Marshal(a)
	if err != nil {
		resp.ErrMsg = cant_create_json
		return
	}
	resp.Data = res
}

type itemResponse struct {
	id          int    `json:"id"`
	name        string `json:"name"`
	description string `json:"description"`
	size        int    `json:"size"`
}

func handleItemsCode(resp *Response, manager *PlayerManager, command Command) {
	a := []itemResponse{}
	for _, elem := range manager.player.Room().AccessibleItems() {
		slt := &itemResponse{
			name:        elem.Name,
			id:          elem.Id,
			description: elem.Description,
			size:        elem.Size,
		}
		a = append(a, *slt)
	}
	res, err := json.Marshal(a)
	if err != nil {
		resp.ErrMsg = cant_create_json
		return
	}
	resp.Data = res
}

func handleBagCode(resp *Response, manager *PlayerManager, command Command) {
	a := []itemResponse{}
	for k, _ := range *manager.player.Bag().Items() {
		it, _ := manager.player.Bag().GetItem(k)
		slt := &itemResponse{
			name:        it.Name,
			id:          it.Id,
			description: it.Description,
			size:        it.Size,
		}
		a = append(a, *slt)
	}
	res, err := json.Marshal(a)
	if err != nil {
		resp.ErrMsg = cant_create_json
		return
	}
	resp.Data = res
}

type intaractiveResponse struct {
	id          int    `json:"id"`
	name        string `json:"name"`
	description string `json:"description"`
}

func handleIteractivesCode(resp *Response, manager *PlayerManager, command Command) {
	a := []intaractiveResponse{}
	for _, elem := range manager.player.Room().Interactives() {
		slt := &intaractiveResponse{
			name:        elem.Name(),
			id:          elem.Id(),
			description: elem.Description(),
		}
		a = append(a, *slt)
	}
	res, err := json.Marshal(a)
	if err != nil {
		resp.ErrMsg = cant_create_json
		return
	}
	resp.Data = res
}

func handleEnterCode(resp *Response, manager *PlayerManager, command Command) {
	var door, ok = manager.player.Room().Doors()[command.itemKey]
	if !ok {
		resp.ErrMsg = fmt.Sprintf(doorNotFoundTemplate, command.itemKey)
		return
	}
	var err = door.Enter(manager.player)
	if err != nil {
		resp.ErrMsg = err.Error()
	}
}

func handleInteractCode(resp *Response, manager *PlayerManager, command Command) {
	var inter, ok = manager.player.Room().Interactives()[command.itemKey]
	if !ok {
		resp.ErrMsg = fmt.Sprintf(interactiveNotFoundTemplate, command.itemKey)
		return
	}
	var result, err = inter.Interact(command.args, command.items)
	if err != nil {
		resp.ErrMsg = err.Error()
		return
	}
	resp.Msg = result.Msg
	resp.Code = result.Code
}

func handleTakeCode(resp *Response, manager *PlayerManager, command Command) {
	if len(command.args) == 0 {
		resp.ErrMsg = itemCodeNotSupplied
		return
	}

	var itemId, parseErr = strconv.Atoi(command.args[0])
	if parseErr != nil {
		resp.ErrMsg = parseErr.Error()
		return
	}

	var moveErr = manager.player.Room().GetItem(itemId, manager.player)
	if moveErr != nil {
		resp.ErrMsg = moveErr.Error()
	}
	return
}

func handlePutCode(resp *Response, manager *PlayerManager, command Command) {
	if len(command.args) == 0 {
		resp.ErrMsg = itemCodeNotSupplied
		return
	}

	var itemId, parseErr = strconv.Atoi(command.args[0])
	if parseErr != nil {
		resp.ErrMsg = parseErr.Error()
		return
	}

	var moveErr = manager.player.Room().PutItem(itemId, manager.player)
	if moveErr != nil {
		resp.ErrMsg = moveErr.Error()
	}
	return
}
