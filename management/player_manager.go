package management

import (
	"errors"
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"strconv"
	"time"
)

const (
	playerInTheVoid     = "Player is not in a room"
	badCode             = "Bad code"
	itemCodeNotSupplied = "Item code not supplied"
	itemsAbsent         = "Not all items are in your bag"

	managerNotStartedCode = iota
	managerWorkCode
	managerFinishedCode
)

func getDoorNotFoundMsg(id int) string {
	return fmt.Sprintf("Door %d not found", id)
}

func getInteractiveNotFoundMsg(id int) string {
	return fmt.Sprintf("Interactive \"%s\" not found", id)
}

func getSlotNotAvailable(id int) string {
	return fmt.Sprintf("Slot %v not found", id)
}

type PlayerManager struct {
	stateCode int
	player    *entities.Player
	inChan    chan Command
	outChan   chan Response
	stopChan  chan interface{}
	timeOut   time.Duration
}

func NewPlayerManager(
	player *entities.Player,
	commandBuffSize int,
	responseBuffSize int,
	timeOut time.Duration,
) (*PlayerManager, error) {
	if player.Room() == nil {
		return nil, errors.New(playerInTheVoid)
	}

	var result = &PlayerManager{
		stateCode: managerNotStartedCode,
		player:    player,
		inChan:    make(chan Command, commandBuffSize),
		outChan:   make(chan Response, responseBuffSize),
		stopChan:  make(chan interface{}, 1),
		timeOut:   timeOut,
	}
	go result.Run()

	return result, nil
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
	var lastUpdate = time.Now()

	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case command := <-manager.inChan:
			var resp = manager.getCommandResponse(command)
			manager.outChan <- resp
			if resp.IsFinish() {
				break
			}
			lastUpdate = time.Now()
		case <-manager.stopChan:
			break
		default:
			if time.Since(lastUpdate) > manager.timeOut {
				break
			}
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
		GetRoomCode:         handleRoomCode,
		GetSlotsCode:        handleSlotsCode,
		GetDoorsCode:        handleDoorsCode,
		GetItemsCode:        handleItemsCode,
		GetBagCode:          handleBagCode,
		GetInteractivesCode: handleInteractivesCode,
		GetSlotContent:      handleSlotContentCode,
		enterCode:           handleEnterCode,
		interactCode:        handleInteractCode,
		takeCode:            handleTakeCode,
		putCode:             handlePutCode,
	}

	var f, ok = methodMap[command.TypeCode]

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
		it, _ := manager.player.Bag().WatchItem(k)
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

func handleSlotContentCode(resp *Response, manager *PlayerManager, command Command) {
	a := []itemResponse{}
	slot, err := manager.player.Room().Slots()[command.ItemKey]
	if err == false {
		resp.ErrMsg = getSlotNotAvailable(command.ItemKey)
		return
	}
	for k := range slot.Items() {
		it, _ := slot.WatchItem(k)
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
	for _, elem := range manager.player.Room().AccessibleInteractives() {
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
	var door, ok = manager.player.Room().Doors()[command.ItemKey]
	if !ok {
		resp.ErrMsg = getDoorNotFoundMsg(command.ItemKey)
		return
	}
	var err = door.Enter(manager.player)
	if err != nil {
		resp.ErrMsg = err.Error()
	}

	resp.Data = "You entered the door"
}

func handleInteractCode(resp *Response, manager *PlayerManager, command Command) {
	var itemsInBag = func(items []entities.Item, bag *entities.Slot) error {
		for _, item := range items {
			var _, ok = bag.Items()[item.Id]
			if !ok {
				return errors.New(itemsAbsent)
			}
		}

		return nil
	}

	var inter, ok = manager.player.Room().AccessibleInteractives()[command.ItemKey]
	if !ok {
		resp.ErrMsg = getInteractiveNotFoundMsg(command.ItemKey)
		return
	}

	var bagErr = itemsInBag(command.Items, manager.Player().Bag())
	if bagErr != nil {
		resp.ErrMsg = bagErr.Error()
		return
	}

	var result, err = inter.Interact(command.Args, command.Items)
	if err != nil {
		resp.ErrMsg = err.Error()
		return
	}
	resp.Msg = result.Msg
	resp.Code = result.Code
}

func handleTakeCode(resp *Response, manager *PlayerManager, command Command) {
	var itemId = command.ItemKey

	if len(command.Args) < 1 {
		resp.ErrMsg = "No slot in command"
		return
	}
	slt := command.Args[0]
	slotID, err := strconv.Atoi(slt)
	if err != nil {
		resp.ErrMsg = "Bad slot id"
		return
	}
	var moveErr = manager.player.Room().GetItem(itemId, slotID, manager.player)
	if moveErr != nil {
		resp.ErrMsg = moveErr.Error()
	}
}

func handlePutCode(resp *Response, manager *PlayerManager, command Command) {
	var itemId = command.ItemKey

	if len(command.Args) < 1 {
		resp.ErrMsg = "No slot in command"
		return
	}
	slt := command.Args[0]
	slotId, err := strconv.Atoi(slt)
	if err != nil {
		resp.ErrMsg = "Bad slot id"
		return
	}
	var moveErr = manager.player.Room().PutItem(itemId, slotId, manager.player)
	if moveErr != nil {
		resp.ErrMsg = moveErr.Error()
	}
}
