package entities

import (
	"errors"
	"fmt"
	"sort"
)

func GetActionNotFoundMsg(id int) string {
	return fmt.Sprintf("Action %d not found", id)
}

func GetActionNotAvailableMsg(id int) string {
	return fmt.Sprintf("Action %d not available", id)
}

func GetFailedToTakeMsg(id int) string {
	return fmt.Sprintf("Failed to take item (id = %d)", id)
}

func GetCanNotPutItemMsg(id int, slotId int) string {
	return fmt.Sprintf("Can not put item \"%d\" to the slot \"%s\"", id, slotId)
}

func GetNoSlotMsg(id int) string {
	return fmt.Sprintf("Can not find slot \"%d\"", id)
}

type SlotsType map[int]*Slot
type DoorsType map[int]*Door
type InteractivesType map[int]InteractiveObject

type Room struct {
	id           int
	name         string
	description  string
	slots        SlotsType
	interactives InteractivesType
	doors        DoorsType
}

func NewRoom(id int, name string, description string) *Room {
	return &Room{
		id:           id,
		name:         name,
		description:  description,
		slots:        make(SlotsType),
		interactives: make(InteractivesType),
		doors:        make(DoorsType),
	}
}

func (r *Room) Id() int {
	return r.id
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

func (r *Room) Slots() SlotsType {
	return r.slots
}

func (r *Room) SetSlots(slots SlotsType) {
	r.slots = slots
}

func (r *Room) AccessibleSlots() (slots []*Slot) {
	return r.accessibleSlots()
}

func (r *Room) GetItem(itemId int, slotId int, player *Player) error {
	slots := r.Slots()
	slot, i := slots[slotId]
	if i == false {
		return errors.New(GetNoSlotMsg (slotId))
	}
	if !slot.IsAccessible() {
		return errors.New(SlotNotAccessible)
	}
	var err = slot.MoveItem(itemId, player.bag)
	if err != nil {
		return errors.New(GetFailedToTakeMsg (itemId))
	}
	return nil
}

func (r *Room) PutItem(itemId int, slotId int, player *Player) error {
	slots := r.Slots()
	slot, i := slots[slotId]
	if i == false {
		return errors.New(GetNoSlotMsg (slotId))
	}
	if !slot.IsAccessible() {
		return errors.New(SlotNotAccessible)
	}
	var err = player.Bag().MoveItem(itemId, slot)
	if err != nil {
		return errors.New(GetCanNotPutItemMsg (itemId, slotId))
	}
	return nil
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

func (r *Room) Interactives() InteractivesType {
	return r.interactives
}

func (r *Room) SetInteractives(interactives InteractivesType) {
	r.interactives = interactives
}

func (r *Room) AccessibleInteractives() InteractivesType {
	var result = make(InteractivesType)
	for _, inter := range r.interactives {
		if inter.IsAccessible() {
			result[inter.Id()] = inter
		}
	}
	return result
}

func (r *Room) Doors() DoorsType {
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
