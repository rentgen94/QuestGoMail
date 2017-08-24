package entities

import "testing"

func TestAction_Name(t *testing.T) {
	var action = NewAction(0, "name", func(l *Labyrinth) (result InteractionResult, err error) {
		return ContinueResult(""), nil
	})

	if action.Name() != "name" {
		t.Fail()
	}
}
