package entities

type Player struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	GameId   int    `json:"game_id"`
	room     *Room
	bag      *Slot
}

func NewPlayer(room *Room, bagCapacity int) *Player {
	return &Player{
		room: room,
		bag:  NewSlot(-1, "bag", bagCapacity, true),
	}
}

func (p *Player) Room() *Room {
	return p.room
}

func (p *Player) SetRoom(room *Room) {
	p.room = room
}

func (p *Player) Bag() *Slot {
	return p.bag
}

func (p *Player) SetBag(bag *Slot) {
	p.bag = bag
}
