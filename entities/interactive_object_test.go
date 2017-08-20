package entities

import (
	"errors"
	"fmt"
	"testing"
)

func getFailRoom() *Room {
	var action = NewAction(
		"",
		false,
		func(r *Room) (res InteractionResult, err error) {
			return ContinueResult(""), nil
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
		true,
		func(r *Room) (res InteractionResult, err error) {
			return ContinueResult("Success"), nil
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
	var testData = []struct {
		room         *Room
		codeGen      actionCodeGeneratorType
		isAccessible bool
		errIsNil     bool
		msg          string
		errMsg       string
	}{
		{
			getSuccessRoom(),
			func(args []string, items []Item) (code int, err error) {
				return 0, nil
			},
			false,
			false,
			"",
			objectNotAccessible,
		},
		{
			getFailRoom(),
			func(args []string, items []Item) (code int, err error) {
				return 0, nil
			},
			true,
			false,
			"",
			fmt.Sprintf(ActionNotAvailableTemplate, 0),
		},
		{
			getSuccessRoom(),
			func(args []string, items []Item) (code int, err error) {
				return 0, errors.New("Error")
			},
			true,
			false,
			"",
			"Error",
		},
		{
			getSuccessRoom(),
			func(args []string, items []Item) (code int, err error) {
				return 0, nil
			},
			true,
			true,
			"Success",
			"",
		},
	}

	for i, item := range testData {
		var inter = NewInteractiveObject(
			"",
			"",
			item.isAccessible,
			item.room,
			item.codeGen,
		)

		var res, err = inter.Interact(nil, nil)

		if (err == nil) != item.errIsNil {
			t.Errorf("Expected (err == nil) = %v, got %v (%d)", item.errIsNil, err == nil, i)
		}

		if err != nil && err.Error() != item.errMsg {
			t.Errorf("Expected errMsg \"%s\", got %s (%d)", item.errMsg, err.Error(), i)
		}

		if res.Msg != item.msg {
			t.Errorf("Expected res \"%s\", got \"%s\" (%d)", item.msg, res, i)
		}
	}
}
