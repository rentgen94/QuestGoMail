package entities

type Player struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	room     *Room
	bag      *Slot
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
