package entities

type Room struct {
	id           int
	doors        []*Door
	slots        []*Slot
	interactives []InteractiveObject
}

func (r *Room) Equals(another *Room) bool {
	return r.id == another.id
}

func (r *Room) AccessibleSlots() []*Slot {
	var result = make([]*Slot, 0)

	for _, slot := range r.slots {
		if slot.IsAccessible() {
			result = append(result, slot)
		}
	}

	return result
}

func (r *Room) AccessibleInteractives() []InteractiveObject {
	var result = make([]InteractiveObject, 0)

	for _, interactiveObject := range r.interactives {
		if interactiveObject.IsAccessible() {
			result = append(result, interactiveObject)
		}
	}

	return result
}

func (r *Room) AccessibleDoors() []*Door {
	var result = make([]*Door, 0)

	for _, door := range r.doors {
		if door.IsAccessible() {
			result = append(result, door)
		}
	}

	return result
}
