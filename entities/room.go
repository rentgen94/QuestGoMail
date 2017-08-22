package entities

import (
	"errors"
	"fmt"
	"sort"
)

const (
	ActionNotFoundTemplate     = "Action %d not found"
	ActionNotAvailableTemplate = "Action %d not available"
	FailedToTakeTemplate       = "Failed to take item (id = %d)"
	CanNotPutItemTemplate      = "Can not put item \"%d\" to the action \"%s\""
)

type slotsType map[int]*Slot
type doorsType map[int]*Door
type interactivesType map[int]InteractiveObject

type Room struct {
	id           int
	name         string
	description  string
	slots        slotsType
	interactives interactivesType
	doors        doorsType
}

func NewRoom(id int, name string, description string) *Room {
	return &Room{
		id:           id,
		name:         name,
		description:  description,
		slots:        make(slotsType),
		interactives: make(interactivesType),
		doors:        make(doorsType),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Description() string {
	return r.description
}

func (r *Room) Equals(another *Room) bool {
	return r.id == another.id
}

func (r *Room) Slots() slotsType {
	return r.slots
}

func (r *Room) AccessibleSlots() (slots []*Slot) {
	return r.accessibleSlots()
}

func (r *Room) GetItem(itemId int, player *Player) error {
	var accessibleSlots = r.accessibleSlots()
	for _, slot := range accessibleSlots {
		var err = slot.MoveItem(itemId, player.bag)
		if err == nil {
			return nil
		}
	}

	return errors.New(fmt.Sprintf(FailedToTakeTemplate, itemId))
}

func (r *Room) PutItem(itemId int, player *Player) error {
	var accessibleSlots = r.accessibleSlots()
	for _, slot := range accessibleSlots {
		var err = player.bag.MoveItem(itemId, slot)
		if err == nil {
			return nil
		}
	}

	return errors.New(fmt.Sprintf(CanNotPutItemTemplate, itemId, r.Name()))
}

func (r *Room) AccessibleItems() (items []Item) {
	var accessibleSlots = r.accessibleSlots()
	for _, slot := range accessibleSlots {
		for _, itemSlice := range slot.items {
			for _, item := range itemSlice {
				items = append(items, item)
			}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items
}

func (r *Room) Interactives() interactivesType {
	return r.interactives
}

func (r *Room) AccessibleInteractives() (interactives []InteractiveObject) {
	for _, inter := range r.interactives {
		if inter.IsAccessible() {
			interactives = append(interactives, inter)
		}
	}
	return
}

func (r *Room) Doors() doorsType {
	return r.doors
}

func (r *Room) AccessibleDoors() (doors []*Door) {
	for _, door := range r.doors {
		if door.IsAccessible() {
			doors = append(doors, door)
		}
	}
	return
}

func (r *Room) accessibleSlots() (slots []*Slot) {
	for _, slot := range r.slots {
		if slot.isAccessible {
			slots = append(slots, slot)
		}
	}
	return
}
