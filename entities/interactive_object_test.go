package entities

import "testing"

func getFailRoom() *Room {
	var action = NewAction(
		"",
		false,
		func(r *Room) (msg string, err error) {
			return "", nil
		},
	)

	var room = NewRoom(
		0,
		"",
		"",
	)

	room.Actions()[0] = action
	return room
}

func getSuccessRoom() *Room {
	var action = NewAction(
		"",
		false,
		func(r *Room) (msg string, err error) {
			return "", nil
		},
	)

	var room = NewRoom(
		0,
		"",
		"",
	)

	room.Actions()[0] = action
	return room
}

func TestBoundInteractiveObject_Interact(t *testing.T) {

}
