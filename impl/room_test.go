package impl

import (
	"fmt"
	"github.com/rentgen94/QuestGoMail/entities"
	"testing"
)

const (
	openCode         = 0
	treasureName     = "Treasure"
	treasureSlotName = "TreasureSlot"
	boxName          = "Box"
)

var onlyRoom = func() *entities.Room {
	var room = entities.NewRoom(
		0,
		"Комната без окон, без дверей",
		"Пустая комната с коробкой посередине",
	)

	room.Slots()[treasureSlotName] = entities.NewSlot(0, treasureSlotName, 1000, false)

	var openAction = entities.NewAction(
		"open",
		true,
		func(r map[int]*entities.Room) (entities.InteractionResult, error) {
			r[0].Slots()[treasureSlotName].SetAccessible(true)
			return entities.ContinueResult("Treasure slot opened"), nil
		},
	)
	room.Actions()[openCode] = openAction

	var box = entities.NewInteractiveObject(
		0,
		boxName,
		"Коробка с сокровищем внутри",
		true,
		room,
		func(args []string, items []entities.Item) (code int, err error) {
			return 0, nil
		},
	)

	room.Interactives()[boxName] = box

	return room
}()

func TestOnlyRoom(t *testing.T) {
	var temp = onlyRoom
	temp.Interactives()[boxName].Interact(nil, nil)
	fmt.Println(temp)
	fmt.Println(onlyRoom)
}
