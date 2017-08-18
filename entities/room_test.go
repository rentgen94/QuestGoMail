package entities

import (
	"errors"
	"fmt"
	"testing"
)

func getRoom(action *Action) *Room {
	var room = NewRoom(0, "room", "")
	room.Actions()[0] = action
	return room
}

func TestRoom_PerformAction(t *testing.T) {
	var testData = []struct {
		action   *Action
		code     int
		msg      string
		isErrNil bool
		errMsg   string
	}{
		{
			action: NewAction("", false, func(r *Room) (msg string, err error) {
				return "", nil
			}),
			code:     0,
			msg:      "",
			isErrNil: false,
			errMsg:   fmt.Sprintf(ActionNotAvailableTemplate, 0),
		},
		{
			action: NewAction("", true, func(r *Room) (msg string, err error) {
				return "", nil
			}),
			code:     1,
			msg:      "",
			isErrNil: false,
			errMsg:   fmt.Sprintf(ActionNotFoundTemplate, 1),
		},
		{
			action: NewAction("", true, func(r *Room) (msg string, err error) {
				return "", errors.New("Error")
			}),
			code:     0,
			msg:      "",
			isErrNil: false,
			errMsg:   "Error",
		},
		{
			action: NewAction("", true, func(r *Room) (msg string, err error) {
				return "Msg", nil
			}),
			code:     0,
			msg:      "Msg",
			isErrNil: true,
			errMsg:   "",
		},
	}

	for i, item := range testData {
		var room = getRoom(item.action)
		var msg, err = room.PerformAction(item.code)

		if msg != item.msg {
			t.Errorf("Expected msg %s, got %s (%d)", item.msg, msg, i)
		}

		if (err == nil) != item.isErrNil {
			t.Errorf("Expected err == nil = %v, got %v (%d)", item.isErrNil, err == nil, i)
		}

		if item.isErrNil && err != nil && err.Error() != item.errMsg {
			t.Errorf("Expected errMsg = \"%s\", got \"%s\" (%d)", item.errMsg, err.Error(), i)
		}
	}
}

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
