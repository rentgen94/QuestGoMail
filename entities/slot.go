package entities

import "errors"

type Slot struct {
	capacity     int
	contains     int
	isAccessible bool
	name         string
	items        map[int][]*Item
}

func NewSlot(name string, capacity int, isAccessible bool) *Slot {
	return &Slot{
		capacity:     capacity,
		contains:     0,
		name:         name,
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

func (s *Slot) canPut(item *Item) error {
	if !s.isAccessible {
		return errors.New("Can slot is not open")
	}
	if item.Size+s.contains > s.capacity {
		return errors.New("Item is too big")
	}
	return nil
}

func (s *Slot) hasItem(id int) error {
	if !s.isAccessible {
		return errors.New("Can slot is not open")
	}

	_, err := s.items[id]

	if !err {
		return errors.New("No such item in this slot")
	}

	return nil
}

func (s *Slot) getItem(id int) (item *Item, err error) {
	err = s.hasItem(id)
	if err != nil {
		return nil, err
	}

	item = s.items[id][len(s.items[id])-1]
	s.items[id] = s.items[id][:len(s.items)-2]
	s.contains -= item.Size
	return item, nil
}

func (s *Slot) putItem(item *Item) error {
	err := s.canPut(item)
	if err != nil {
		return err
	}

	s.contains += item.Size
	s.items[item.Id] = append(s.items[item.Id], item)

	return nil
}
