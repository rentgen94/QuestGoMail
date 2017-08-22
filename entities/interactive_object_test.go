package entities

import (
	"errors"
	"fmt"
	"testing"
)

func getFailAction() *Action {
	return NewAction(
		"",
		nil,
		false,
		func(l *Labyrinth) (res InteractionResult, err error) {
			return ContinueResult(""), nil
		},
	)
}

func getSuccessAction() *Action {
	return NewAction(
		"",
		nil,
		true,
		func(l *Labyrinth) (res InteractionResult, err error) {
			return ContinueResult("Success"), nil
		},
	)
}

func TestBoundInteractiveObject_Interact(t *testing.T) {
	var testData = []struct {
		action       *Action
		checker InputChecker
		isAccessible bool
		errIsNil     bool
		msg          string
		errMsg       string
	}{
		{
			getSuccessAction(),
			func(args []string, items []Item) error {
				return nil
			},
			false,
			false,
			"",
			objectNotAccessible,
		},
		{
			getFailAction(),
			func(args []string, items []Item) error {
				return nil
			},
			true,
			false,
			"",
			fmt.Sprintf(ActionNotAvailableTemplate, 0),
		},
		{
			getSuccessAction(),
			func(args []string, items []Item) error {
				return errors.New("Error")
			},
			true,
			false,
			"",
			"Error",
		},
		{
			getSuccessAction(),
			func(args []string, items []Item) error {
				return nil
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
			0,
			"",
			item.isAccessible,
			item.checker,
			item.action,
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
