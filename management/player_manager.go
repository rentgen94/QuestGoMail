package management

import (
	"errors"
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"strconv"
)

const (
	playerInTheVoid             = "Player is not in a room"
	badCode                     = "Bad code"
	doorNotFoundTemplate        = "Door %v not found"
	interactiveNotFoundTemplate = "Interactive \"%s\" not found"
	itemCodeNotSupplied         = "Item code not supplied"

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

func (manager *PlayerManager) Player() *entities.Player {
	return manager.player
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
		GetIteractivesCode: handleInteractivesCode,
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

func handleRoomCode(resp *Response, manager *PlayerManager, command Command) {
	type roomResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	a := &roomResponse{
		Name:        manager.player.Room().Name(),
		Description: manager.player.Room().Description(),
	}
	resp.Data = a
}

func handleSlotsCode(resp *Response, manager *PlayerManager, command Command) {
	type slotsResponse struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Capacity int    `json:"capacity"`
		Contains int    `json:"contains"`
	}

	a := []slotsResponse{}
	for _, elem := range manager.player.Room().Slots() {
		slt := &slotsResponse{
			Name:     elem.Name(),
			Capacity: elem.Capacity(),
			Contains: elem.Contains(),
			Id:       elem.Id(),
		}
		a = append(a, *slt)
	}
	resp.Data = a
}

func handleDoorsCode(resp *Response, manager *PlayerManager, command Command) {
	type doorsResponse struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	a := []doorsResponse{}
	for _, elem := range manager.player.Room().Doors() {
		slt := &doorsResponse{
			Name: elem.Name(),
			Id:   elem.Id(),
		}
		a = append(a, *slt)
	}
	resp.Data = a
}

type itemResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
}

func handleItemsCode(resp *Response, manager *PlayerManager, command Command) {
	a := []itemResponse{}
	for _, elem := range manager.player.Room().AccessibleItems() {
		slt := &itemResponse{
			Name:        elem.Name,
			Id:          elem.Id,
			Description: elem.Description,
			Size:        elem.Size,
		}
		a = append(a, *slt)
	}
	resp.Data = a
}

func handleBagCode(resp *Response, manager *PlayerManager, command Command) {
	a := []itemResponse{}
	for k := range manager.player.Bag().Items() {
		it, _ := manager.player.Bag().GetItem(k)
		slt := &itemResponse{
			Name:        it.Name,
			Id:          it.Id,
			Description: it.Description,
			Size:        it.Size,
		}
		a = append(a, *slt)
	}
	resp.Data = a
}

func handleInteractivesCode(resp *Response, manager *PlayerManager, command Command) {
	type interactiveResponse struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	a := []interactiveResponse{}
	for _, elem := range manager.player.Room().Interactives() {
		slt := &interactiveResponse{
			Name:        elem.Name(),
			Id:          elem.Id(),
			Description: elem.Description(),
		}
		a = append(a, *slt)
	}
	resp.Data = a
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
}
