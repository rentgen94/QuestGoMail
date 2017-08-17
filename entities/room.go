package entities

import (
	"errors"
	"fmt"
)

const (
	actionNotFoundTemplate = "Action %d not found"
	actionNotAvailableTemplate = "Action %d not available"
)

type slotsType map[string]*Slot
type doorsType map[string]*Door
type interactivesType map[string]InteractiveObject
type actionsType map[int]*Action

type Room struct {
	id           int
	name         string
	description  string
	slots        slotsType
	interactives interactivesType
	doors        doorsType
	actions      actionsType
}

func NewRoom(id int, name string, description string) *Room {
	return &Room{
		id:           id,
		name:         name,
		description:  description,
		slots:        make(slotsType),
		interactives: make(interactivesType),
		doors:        make(doorsType),
		actions:      make(actionsType),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Description() string {
	return r.description
}

func (r *Room) PerformAction(actionCode int) (string, error) {
	var action, ok = r.actions[actionCode]
	if !ok {
		return "", errors.New(fmt.Sprintf(actionNotFoundTemplate, actionCode))
	}

	if !action.isAccessible {
		return "", errors.New(fmt.Sprintf(actionNotAvailableTemplate, actionCode))
	}

	return action.act(r)
}

func (r *Room) Equals(another *Room) bool {
	return r.id == another.id
}

func (r *Room) Slots() slotsType {
	return r.slots
}

func (r *Room) Interactives() interactivesType {
	return r.interactives
}

func (r *Room) Doors() doorsType {
	return r.doors
}

func (r *Room) Actions() actionsType {
	return r.actions
}
