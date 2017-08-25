package entities

type Player struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	GameId    int    `json:"game_id"`
	room      *Room
	labyrinth *Labyrinth
	bag       *Slot
}

func NewPlayer(lab *Labyrinth, bagCapacity int) *Player {
	return &Player{
		room:      lab.StartRoom(),
		labyrinth: lab,
		bag:       NewSlot(-1, "bag", bagCapacity, true),
	}
}

func (p *Player) Equal(another *Player) bool {
	if p == nil || another == nil {
		return false
	} else if p.Login != another.Login {
		return false
	} else if p.Password != another.Password {
		return false
	}

	return true
}

func (p *Player) Labyrinth() *Labyrinth {
	return p.labyrinth
}

func (p *Player) SetLabyrinth(lab *Labyrinth) {
	p.labyrinth = lab
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
