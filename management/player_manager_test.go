package management

import (
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"strconv"
	"testing"
)

const (
	interactiveName = "Box"
	slotName        = "Slot"
	doorName        = "Only"
	itemName        = "Item"
	itemId          = 100
	bigItemId       = 200
	playerItemId    = 300
	doorId          = 10
	interactiveId   = 15
)

func getRoom() *entities.Room {
	var room = entities.NewRoom(0, "", "")
	var action = entities.NewAction(
		"open",
		func(labyrinth *entities.Labyrinth) (result entities.InteractionResult, err error) {
			return entities.ContinueResult("Success"), nil
		},
	)

	var box = entities.NewInteractiveObject(
		interactiveId,
		interactiveName,
		"",
		true,
		func(args []string, items []entities.Item) error {
			return nil
		},
		action,
	)

	room.Interactives()[box.Id()] = box

	var slot = entities.NewSlot(0, slotName, 10, true)
	room.Slots()[slot.Id()] = slot

	var item = entities.Item{Name: itemName, Size: 1, Id: itemId}
	slot.PutItem(item)

	var bigItem = entities.Item{Name: itemName, Size: 6, Id: bigItemId}
	slot.PutItem(bigItem)

	var door = entities.NewDoor(doorId, doorName, true, room, room)
	room.Doors()[door.Id()] = door

	return room
}

func getPlayer() *entities.Player {
	var player = entities.NewPlayer(getRoom(), 5)
	player.Bag().PutItem(entities.Item{Id: playerItemId, Size: 0})
	return player
}

func TestPlayerManager_getCommandResponse(t *testing.T) {
	var testData = []struct {
		command Command
		errMsg  string
	}{
		{
			command: NewCommand(GetRoomCode, 0, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(GetBagCode, 0, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, doorId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, 1, nil, nil),
			errMsg:  fmt.Sprintf(doorNotFoundTemplate, 1),
		},
		{
			command: NewCommand(interactCode, interactiveId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(interactCode, 100, nil, nil),
			errMsg:  fmt.Sprintf(interactiveNotFoundTemplate, 100),
		},
		{
			command: NewCommand(takeCode, 0, make([]string, 0), nil),
			errMsg:  itemCodeNotSupplied,
		},
		{
			command: NewCommand(takeCode, 0, []string{strconv.Itoa(itemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(takeCode, 0, []string{strconv.Itoa(itemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, itemId+1),
		},
		{
			command: NewCommand(takeCode, 0, []string{strconv.Itoa(bigItemId)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, bigItemId),
		},
		{
			command: NewCommand(putCode, 0, []string{strconv.Itoa(playerItemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(putCode, 0, []string{strconv.Itoa(playerItemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.CanNotPutItemTemplate, playerItemId+1, ""),
		},
		{
			command: NewCommand(-100, 0, nil, nil),
			errMsg:  badCode,
		},
	}

	for i, item := range testData {
		var player = getPlayer()
		var manager, err = NewPlayerManager(player, 10, 10)
		if err != nil {
			t.Fatal(err)
		}
		go manager.Run()
		manager.CommandChan() <- item.command
		var resp = <-manager.RespChan()
		if resp.ErrMsg != item.errMsg {
			t.Errorf("Expected ErrMsg \"%s\", got \"%f\" (%v)", item.errMsg, resp.ErrMsg, i)
		}
		manager.Stop()
	}
}

func TestPlayerManager_Run(t *testing.T) {
	var testData = []struct {
		command Command
		errMsg  string
	}{
		{
			command: NewCommand(GetRoomCode, 0, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(GetBagCode, 0, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, doorId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, 100, nil, nil),
			errMsg:  fmt.Sprintf(doorNotFoundTemplate, 100),
		},
		{
			command: NewCommand(interactCode, interactiveId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(interactCode, 100, nil, nil),
			errMsg:  fmt.Sprintf(interactiveNotFoundTemplate, 100),
		},
		{
			command: NewCommand(takeCode, 0, make([]string, 0), nil),
			errMsg:  itemCodeNotSupplied,
		},
		{
			command: NewCommand(takeCode, 0, []string{strconv.Itoa(itemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(takeCode, 0, []string{strconv.Itoa(itemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, itemId+1),
		},
		{
			command: NewCommand(takeCode, 0, []string{strconv.Itoa(bigItemId)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, bigItemId),
		},
		{
			command: NewCommand(putCode, 0, []string{strconv.Itoa(playerItemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(putCode, 0, []string{strconv.Itoa(playerItemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.CanNotPutItemTemplate, playerItemId+1, ""),
		},
		{
			command: NewCommand(-100, 0, nil, nil),
			errMsg:  badCode,
		},
	}

	for i, item := range testData {
		var player = getPlayer()
		var manager, err = NewPlayerManager(player, 10, 10)
		if err != nil {
			t.Fatal(err)
		}
		go manager.Run()
		manager.inChan <- item.command
		var resp = <-manager.outChan

		if resp.ErrMsg != item.errMsg {
			t.Errorf("Expected ErrMsg \"%s\", got \"%f\" (%v)", item.errMsg, resp.ErrMsg, i)
		}
		manager.Stop()
	}
}
