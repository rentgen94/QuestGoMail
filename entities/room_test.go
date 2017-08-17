package entities

import (
	"testing"
	"fmt"
	"errors"
)

func getRoom(action *Action) *Room {
	var room = NewRoom(0, "room", "")
	room.Actions()[0] = action
	return room
}

func TestRoom_PerformAction(t *testing.T) {
	var testData = []struct {
		action *Action
		code   int
		msg    string
		isErrNil bool
		errMsg string
	}{
		{
			action:NewAction("", false, func(r *Room) (msg string, err error) {
				return "", nil
			}),
			code:   0,
			msg:    "",
			isErrNil: false,
			errMsg: fmt.Sprintf(actionNotAvailableTemplate, 0),
		},
		{
			action:NewAction("", true, func(r *Room) (msg string, err error) {
				return "", nil
			}),
			code:   1,
			msg:    "",
			isErrNil: false,
			errMsg: fmt.Sprintf(actionNotFoundTemplate, 1),
		},
		{
			action:NewAction("", true, func(r *Room) (msg string, err error) {
				return "", errors.New("Error")
			}),
			code:   0,
			msg:    "",
			isErrNil: false,
			errMsg: "Error",
		},
		{
			action:NewAction("", true, func(r *Room) (msg string, err error) {
				return "Msg", nil
			}),
			code:   0,
			msg:    "Msg",
			isErrNil: true,
			errMsg: "",
		},
	}

	for i, item := range testData {
		var room = getRoom(item.action)
		var msg, err = room.PerformAction(item.code)

		if msg != item.msg {
			t.Error(fmt.Sprintf("Expected msg %s, got %s (%d)", item.msg, msg, i))
		}

		if (err == nil) != item.isErrNil {
			t.Error(fmt.Sprintf("Expected err == nil = %v, got %v (%d)", item.isErrNil, err == nil, i))
		}

		if item.isErrNil && err != nil && err.Error() != item.errMsg {
			t.Error(fmt.Sprintf("Expected errMsg = \"%s\", got \"%s\" (%d)", item.errMsg, err.Error(), i))
		}
	}
}
