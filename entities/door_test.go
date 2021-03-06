package entities

import (
	"testing"
)

func TestDoor_Enter(t *testing.T) {
	var room1 = NewRoom(0, "Room1", "")
	var room2 = NewRoom(1, "Room2", "")
	var room3 = NewRoom(2, "Room3", "")

	var door = NewDoor(10, "", true)
	door.SetRoom1(room1)
	door.SetRoom2(room2)

	var player = &Player{room: room1}

	if player.room.name != "Room1" {
		t.Errorf("Expected %s, got %s", "Room1", player.room.name)
	}

	door.Enter(player)

	if player.room.name != "Room2" {
		t.Errorf("Expected %s, got %s", "Room2", player.room.name)
	}

	player.room = room3

	var err = door.Enter(player)

	if err == nil {
		t.Error("Had to fail")
	}
}
