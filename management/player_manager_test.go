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
)

func getRoom() *entities.Room {
	var room = entities.NewRoom(0, "", "")
	var action = entities.NewAction(
		"open",
		true,
		func(r *entities.Room) (msg string, err error) {
			return "Success", nil
		},
	)
	room.Actions()[0] = action

	var box = entities.NewInteractiveObject(
		interactiveName,
		"",
		true,
		room,
		func(args []string, items []entities.Item) (code int, err error) {
			return 0, nil
		},
	)
	room.Interactives()[box.Name()] = box

	var slot = entities.NewSlot(slotName, 10, true)
	room.Slots()[slot.Name()] = slot

	var item = entities.Item{Name: itemName, Size: 1, Id: itemId}
	slot.PutItem(item)

	var bigItem = entities.Item{Name: itemName, Size: 6, Id: bigItemId}
	slot.PutItem(bigItem)

	var door = entities.NewDoor(doorName, true, room, room)
	room.Doors()[door.Name()] = door

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
			command: NewCommand(getSlotsCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(getItemsCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(getInteractivesCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(getDoorsCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, doorName, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, "Some", nil, nil),
			errMsg:  fmt.Sprintf(doorNotFoundTemplate, "Some"),
		},
		{
			command: NewCommand(interactCode, interactiveName, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(interactCode, "Some", nil, nil),
			errMsg:  fmt.Sprintf(interactiveNotFoundTemplate, "Some"),
		},
		{
			command: NewCommand(takeCode, "", make([]string, 0), nil),
			errMsg:  itemCodeNotSupplied,
		},
		{
			command: NewCommand(takeCode, "", []string{strconv.Itoa(itemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(takeCode, "", []string{strconv.Itoa(itemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, itemId+1),
		},
		{
			command: NewCommand(takeCode, "", []string{strconv.Itoa(bigItemId)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, bigItemId),
		},
		{
			command: NewCommand(putCode, "", []string{strconv.Itoa(playerItemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(putCode, "", []string{strconv.Itoa(playerItemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.CanNotPutItemTemplate, playerItemId+1, ""),
		},
		{
			command: NewCommand(-100, "", nil, nil),
			errMsg:  badCode,
		},
	}

	for i, item := range testData {
		var player = getPlayer()
		var manager, err = NewPlayerManager(player)
		if err != nil {
			t.Fatal(err)
		}

		var resp = manager.getCommandResponse(item.command)
		if resp.errMsg != item.errMsg {
			t.Errorf("Expected errMsg \"%s\", got \"%f\" (%v)", item.errMsg, resp.errMsg, i)
		}
	}
}

func TestPlayerManager_Run(t *testing.T) {
	var testData = []struct {
		command Command
		errMsg  string
	}{
		{
			command: NewCommand(getSlotsCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(getItemsCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(getInteractivesCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(getDoorsCode, "", nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, doorName, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(enterCode, "Some", nil, nil),
			errMsg:  fmt.Sprintf(doorNotFoundTemplate, "Some"),
		},
		{
			command: NewCommand(interactCode, interactiveName, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(interactCode, "Some", nil, nil),
			errMsg:  fmt.Sprintf(interactiveNotFoundTemplate, "Some"),
		},
		{
			command: NewCommand(takeCode, "", make([]string, 0), nil),
			errMsg:  itemCodeNotSupplied,
		},
		{
			command: NewCommand(takeCode, "", []string{strconv.Itoa(itemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(takeCode, "", []string{strconv.Itoa(itemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, itemId+1),
		},
		{
			command: NewCommand(takeCode, "", []string{strconv.Itoa(bigItemId)}, nil),
			errMsg:  fmt.Sprintf(entities.FailedToTakeTemplate, bigItemId),
		},
		{
			command: NewCommand(putCode, "", []string{strconv.Itoa(playerItemId)}, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(putCode, "", []string{strconv.Itoa(playerItemId + 1)}, nil),
			errMsg:  fmt.Sprintf(entities.CanNotPutItemTemplate, playerItemId+1, ""),
		},
		{
			command: NewCommand(-100, "", nil, nil),
			errMsg:  badCode,
		},
	}

	for i, item := range testData {
		var player = getPlayer()
		var manager, err = NewPlayerManager(player)
		if err != nil {
			t.Fatal(err)
		}
		go manager.Run()
		manager.inChan <- item.command
		var resp = <-manager.outChan

		if resp.errMsg != item.errMsg {
			t.Errorf("Expected errMsg \"%s\", got \"%f\" (%v)", item.errMsg, resp.errMsg, i)
		}
	}
}
