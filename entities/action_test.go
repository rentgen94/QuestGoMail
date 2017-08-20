package entities

import "testing"

func TestAction_IsAccessible(t *testing.T) {
	var action = NewAction("name", true, func(r *Room) (result InteractionResult, err error) {
		return ContinueResult(""), nil
	})

	if !action.IsAccessible() {
		t.Fail()
	}

	action.SetAccessible(false)

	if action.IsAccessible() {
		t.Fail()
	}
}

func TestAction_Name(t *testing.T) {
	var action = NewAction("name", true, func(r *Room) (result InteractionResult, err error) {
		return ContinueResult(""), nil
	})

	if action.Name() != "name" {
		t.Fail()
	}
}
