package management

import (
	"github.com/rentgen94/QuestGoMail/entities"
	"testing"
	"time"
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
		0,
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

	var door = entities.NewDoor(doorId, doorName, true)
	door.SetRoom1(room)
	door.SetRoom2(room)
	room.Doors()[door.Id()] = door

	return room
}

func getPlayer() *entities.Player {
	var lab = entities.NewLabyrinth(0, "", "")
	lab.SetStartRoom(getRoom())
	var player = entities.NewPlayer(lab, 5)
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
			errMsg:  getDoorNotFoundMsg(1),
		},
		{
			command: NewCommand(interactCode, interactiveId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(interactCode, 100, nil, nil),
			errMsg:  getInteractiveNotFoundMsg(100),
		},
		{
			command: NewCommand(takeCode, itemId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(takeCode, itemId+1, nil, nil),
			errMsg:  entities.GetFailedToTakeMsg(itemId + 1),
		},
		{
			command: NewCommand(takeCode, bigItemId, nil, nil),
			errMsg:  entities.GetFailedToTakeMsg(bigItemId),
		},
		{
			command: NewCommand(putCode, playerItemId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(putCode, playerItemId+1, nil, nil),
			errMsg:  entities.GetCanNotPutItemMsg(playerItemId+1, ""),
		},
		{
			command: NewCommand(-100, 0, nil, nil),
			errMsg:  badCode,
		},
	}

	for i, item := range testData {
		var player = getPlayer()
		var manager, err = NewPlayerManager(player, 10, 10, time.Hour)
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
			errMsg:  getDoorNotFoundMsg(100),
		},
		{
			command: NewCommand(interactCode, interactiveId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(interactCode, 100, nil, nil),
			errMsg:  getInteractiveNotFoundMsg(100),
		},
		{
			command: NewCommand(takeCode, itemId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(takeCode, itemId+1, nil, nil),
			errMsg:  entities.GetFailedToTakeMsg(itemId + 1),
		},
		{
			command: NewCommand(takeCode, bigItemId, nil, nil),
			errMsg:  entities.GetFailedToTakeMsg(bigItemId),
		},
		{
			command: NewCommand(putCode, playerItemId, nil, nil),
			errMsg:  "",
		},
		{
			command: NewCommand(putCode, playerItemId+1, nil, nil),
			errMsg:  entities.GetCanNotPutItemMsg(playerItemId+1, ""),
		},
		{
			command: NewCommand(-100, 0, nil, nil),
			errMsg:  badCode,
		},
	}

	for i, item := range testData {
		var player = getPlayer()
		var manager, err = NewPlayerManager(player, 10, 10, time.Hour)
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
