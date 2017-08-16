package entities

import "sort"

type Bag struct {
	slots    []*Slot
}

func NewBag(slots []*Slot) *Bag {
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].capacity < slots[j].capacity
	})	// TODO maybe should copy slice before sorting
	return &Bag{
		slots:slots,
	}
}

// Альтернатива - просто завести вместимость сумки и удалить слоты
type AltBag struct {
	capacity int
	contains int
	items map[string]*Item
}


