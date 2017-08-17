package entities

import "testing"

func TestAction_IsAccessible(t *testing.T) {
	var action = NewAction("name", true, func(r *Room) (msg string, err error) {
		return "", nil
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
	var action = NewAction("name", true, func(r *Room) (msg string, err error) {
		return "", nil
	})

	if action.Name() != "name" {
		t.Fail()
	}
}
