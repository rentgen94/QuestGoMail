package entities

import (
	"testing"
)

func TestRoom_Equals(t *testing.T) {
	var testData = []struct {
		room1   *Room
		room2   *Room
		isEqual bool
	}{
		{
			NewRoom(0, "Some", "d"),
			NewRoom(0, "Any", "I"),
			true,
		},
		{
			NewRoom(0, "", ""),
			NewRoom(1, "", ""),
			false,
		},
	}

	for i, item := range testData {
		if item.room1.Equals(item.room2) != item.isEqual {
			t.Errorf("Expected %v, got %v (%d)", item.isEqual, item.room1.Equals(item.room2), i)
		}
	}
}
