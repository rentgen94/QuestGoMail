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
	doorNotFoundTemplate        = "Door \"%s\" not found"
	interactiveNotFoundTemplate = "Interactive \"%s\" not found"
	itemCodeNotSupplied         = "Item code not supplied"
)

type PlayerManager struct {
	player     *entities.Player
	inChan     chan Command
	outChan    chan Response
	finishChan chan interface{}
}

func NewPlayerManager(player *entities.Player) (*PlayerManager, error) {
	if player.Room() == nil {
		return nil, errors.New(playerInTheVoid)
	}

	return &PlayerManager{
		player:  player,
		inChan:  make(chan Command),
		outChan: make(chan Response),
	}, nil
}

func (manager *PlayerManager) Run() {
	// TODO add possibility to break by victory
	for {
		select {
		case command := <-manager.inChan:
			var resp = manager.getCommandResponse(command)
			manager.outChan <- resp
		case <-manager.finishChan:
			break
		}
	}

}

func (manager *PlayerManager) getCommandResponse(command Command) (resp Response) {
	var methodMap = map[int]func(*Response, *PlayerManager, Command){
		getDoorsCode:        handleDoorsCode,
		getSlotsCode:        handleSlotsCode,
		getInteractivesCode: handleInteractivesCode,
		getItemsCode:        handleItemsCode,
		enterCode:           handleEnterCode,
		interactCode:        handleInteractCode,
		takeCode:            handleTakeCode,
		putCode:             handlePutCode,
	}

	var f, ok = methodMap[command.typeCode]
	if !ok {
		resp.errMsg = badCode
		return
	}
	f(&resp, manager, command)
	return
}

func handleDoorsCode(resp *Response, manager *PlayerManager, command Command) {
	resp.data = manager.player.Room().AccessibleDoors()
}

func handleSlotsCode(resp *Response, manager *PlayerManager, command Command) {
	resp.data = manager.player.Room().AccessibleSlots()
}

func handleInteractivesCode(resp *Response, manager *PlayerManager, command Command) {
	resp.data = manager.player.Room().AccessibleInteractives()
}

func handleItemsCode(resp *Response, manager *PlayerManager, command Command) {
	resp.data = manager.player.Room().AccessibleItems()
}

func handleEnterCode(resp *Response, manager *PlayerManager, command Command) {
	var door, ok = manager.player.Room().Doors()[command.itemKey]
	if !ok {
		resp.errMsg = fmt.Sprintf(doorNotFoundTemplate, command.itemKey)
		return
	}
	var err = door.Enter(manager.player)
	if err != nil {
		resp.errMsg = err.Error()
	}
}

func handleInteractCode(resp *Response, manager *PlayerManager, command Command) {
	var inter, ok = manager.player.Room().Interactives()[command.itemKey]
	if !ok {
		resp.errMsg = fmt.Sprintf(interactiveNotFoundTemplate, command.itemKey)
		return
	}
	var msg, err = inter.Interact(command.args, command.items)
	if err != nil {
		resp.errMsg = err.Error()
		return
	}
	resp.msg = msg
}

func handleTakeCode(resp *Response, manager *PlayerManager, command Command) {
	if len(command.args) == 0 {
		resp.errMsg = itemCodeNotSupplied
		return
	}

	var itemId, parseErr = strconv.Atoi(command.args[0])
	if parseErr != nil {
		resp.errMsg = parseErr.Error()
		return
	}

	var moveErr = manager.player.Room().GetItem(itemId, manager.player)
	if moveErr != nil {
		resp.errMsg = moveErr.Error()
	}
	return
}

func handlePutCode(resp *Response, manager *PlayerManager, command Command) {
	if len(command.args) == 0 {
		resp.errMsg = itemCodeNotSupplied
		return
	}

	var itemId, parseErr = strconv.Atoi(command.args[0])
	if parseErr != nil {
		resp.errMsg = parseErr.Error()
		return
	}

	var moveErr = manager.player.Room().PutItem(itemId, manager.player)
	if moveErr != nil {
		resp.errMsg = moveErr.Error()
	}
	return
}
