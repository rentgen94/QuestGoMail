package entities

import (
	"errors"
	"fmt"
)

const (
	SlotNotAccessible      = "Slot is not accessible"
	ItemTooBig             = "Item is too big"
	ItemNotPresentTemplate = "Item id = %d not present"
)

func getItemNotPresentMsg(id int) string {
	return fmt.Sprintf("Item id = %d not present", id)
}

type itemsType map[int][]Item

type Slot struct {
	id           int
	capacity     int
	contains     int
	isAccessible bool
	name         string
	items        itemsType
}

func NewSlot(id int, name string, capacity int, isAccessible bool) *Slot {
	return &Slot{
		id:           id,
		capacity:     capacity,
		contains:     0,
		name:         name,
		isAccessible: isAccessible,
		items:        make(itemsType),
	}
}

func (s *Slot) Capacity() int {
	return s.capacity
}

func (s *Slot) Id() int {
	return s.id
}

func (s *Slot) Contains() int {
	return s.contains
}

func (s *Slot) Items() itemsType {
	return s.items
}

func (s *Slot) IsAccessible() bool {
	return s.isAccessible
}

func (s *Slot) SetAccessible(isAccessible bool) {
	s.isAccessible = isAccessible
}

func (s *Slot) Name() string {
	return s.name
}

func (s *Slot) TakeItem(id int) (item Item, err error) {
	return s.takeItem(id)
}

func (s *Slot) PutItem(item Item) error {
	return s.putItem(item)
}

func (s *Slot) WatchItem(id int) (Item, error) {
	return s.watchItem(id)
}

func (s *Slot) MoveItem(itemId int, another *Slot) error {
	// TODO synchronize on this method
	var item, getErr = s.takeItem(itemId)
	if getErr != nil {
		return getErr
	}

	var putErr = another.putItem(item)
	if putErr != nil {
		s.putItem(item) // if not sync, may fail here
		return putErr
	}

	return nil
}

func (s *Slot) watchItem(id int) (item Item, err error) {
	var has = s.hasItem(id)
	if !has {
		return Item{}, errors.New(getItemNotPresentMsg(id))
	}

	return s.items[id][0], nil
}

func (s *Slot) takeItem(id int) (item Item, err error) {
	var has = s.hasItem(id)
	if !has {
		return Item{}, errors.New(getItemNotPresentMsg(id))
	}

	item = s.items[id][len(s.items[id])-1]
	s.items[id] = s.items[id][:len(s.items[id])-1]
	if len(s.items[id]) == 0 {
		delete(s.items, id)
	}

	s.contains -= item.Size
	return item, nil
}

func (s *Slot) putItem(item Item) error {
	err := s.canPut(item)
	if err != nil {
		return err
	}

	s.contains += item.Size
	s.items[item.Id] = append(s.items[item.Id], item)

	return nil
}

func (s *Slot) canPut(item Item) error {
	//if !s.isAccessible {	commented out cos while constructing labyrinth it may be necessary to put item to the inaccessible slot1
	//	return errors.New(SlotNotAccessible)
	//}
	if item.Size+s.contains > s.capacity {
		return errors.New(ItemTooBig)
	}
	return nil
}

func (s *Slot) hasItem(id int) bool {
	_, has := s.items[id]
	return has
}
