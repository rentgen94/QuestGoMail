package entities

import "errors"

type Slot struct {
	name         string
	capacity     int
	isAccessible bool
	item         *Item
}

func NewSlot(name string, capacity int, isAccessible bool) *Slot {
	return &Slot{
		name:         name,
		capacity:     capacity,
		item:         nil,
		isAccessible: isAccessible,
	}
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

func (s *Slot) Empty() bool {
	return s.item == nil
}

func (s *Slot) MoveTo(another *Slot) error {
	if s.Empty() {
		return errors.New("Origin slot is empty")
	}
	var putErr = another.canPut(s.item)
	if putErr != nil {
		return putErr
	}

	another.item = s.item
	s.item = nil
	return nil
}

func (s *Slot) canPut(item *Item) error {
	if !s.isAccessible {
		return errors.New("Can slot is not open")
	}
	if item.Size > s.capacity {
		return errors.New("Item is too big")
	}
	return nil
}
