package entities

import "errors"

type Door struct {
	isAccessible bool
	room1        *Room
	room2        *Room
}

func (door *Door) IsAccessible() bool {
	return door.isAccessible
}

func (door *Door) SetAccessible(isAccessible bool) {
	door.isAccessible = isAccessible
}

func (door *Door) Enter(player *Player) error {
	if !player.room.Equals(door.room1) && !player.room.Equals(door.room2) {
		return errors.New("Player is not in the neighbour room")
	}
	if !door.isAccessible {
		return errors.New("Door is closed")
	}

	if player.room.Equals(door.room1) {
		player.room = door.room2
	} else {
		player.room = door.room1
	}
	return nil
}
