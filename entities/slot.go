package entities

import "errors"

type Slot struct {
	capacity      int
	isAccessible  bool
	itemValidator func(item *Item) bool
	item          *Item
}

func DefaultSlot(capacity int) *Slot {
	return &Slot{
		capacity:      capacity,
		itemValidator: func(item *Item) bool { return true },
		item:          nil,
	}
}

func (s *Slot) IsAccessible() bool {
	return s.isAccessible
}

func (s *Slot) SetAccessible(isAccessible bool) {
	s.isAccessible = isAccessible
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
	if !s.itemValidator(item) {
		return errors.New("Item does not fit the slot")
	}
	if item.size > s.capacity {
		return errors.New("Item is too big")
	}
	return nil
}
