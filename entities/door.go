package entities

import "errors"

const (
	wrongPlayerNum = "Player is not in the neighbour action"
	wrongRoom      = "Room is not neighbour to the door"
	doorClosed     = "Door is closed"
)

type Door struct {
	id           int
	name         string
	isAccessible bool
	room1        *Room
	room2        *Room
}

func NewDoor(id int, name string, isAccessible bool) *Door {
	return &Door{
		id:           id,
		name:         name,
		isAccessible: isAccessible,
		room1:        nil,
		room2:        nil,
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

func (door *Door) Room1() *Room {
	return door.room1
}

func (door *Door) Room2() *Room {
	return door.room1
}

func (door *Door) AnotherRoom(room *Room) (*Room, error) {
	if room.Equals(door.room1) {
		return door.room2, nil
	}

	if room.Equals(door.room2) {
		return door.room1, nil
	}

	return nil, errors.New(wrongRoom)
}

func (door *Door) SetRoom1(room *Room) {
	door.room1 = room
}

func (door *Door) SetRoom2(room *Room) {
	door.room2 = room
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
