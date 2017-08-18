package entities

type Player struct {
	room *Room
	bag  *Slot
}

func NewPlayer(room *Room, bagCapacity int) *Player {
	return &Player{
		room: room,
		bag:  NewSlot("bag", bagCapacity, true),
	}
}

func (p *Player) Room() *Room {
	return p.room
}

func (p *Player) Bag() *Slot {
	return p.bag
}
