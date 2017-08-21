package entities

import "errors"

const (
	wrongPlayerNum = "Player is not in the neighbour room"
	doorClosed     = "Door is closed"
)

type Door struct {
	id           int
	name         string
	isAccessible bool
	room1        *Room
	room2        *Room
}

func NewDoor(name string, isAccessible bool, room1 *Room, room2 *Room) *Door {
	return &Door{
		name:         name,
		isAccessible: isAccessible,
		room1:        room1,
		room2:        room2,
	}
}

func (door *Door) IsAccessible() bool {
	return door.isAccessible
}

func (door *Door) SetAccessible(isAccessible bool) {
	door.isAccessible = isAccessible
}

func (door *Door) Name() string {
	return door.name
}

func (door *Door) Id() int {
	return door.id
}

func (door *Door) Enter(player *Player) error {
	if !player.room.Equals(door.room1) && !player.room.Equals(door.room2) {
		return errors.New(wrongPlayerNum)
	}
	if !door.isAccessible {
		return errors.New(doorClosed)
	}

	if player.room.Equals(door.room1) {
		player.room = door.room2
	} else {
		player.room = door.room1
	}
	return nil
}
