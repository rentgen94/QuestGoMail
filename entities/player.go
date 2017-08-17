package entities

type Player struct {
	room *Room
	// TODO add content
}

func (p *Player) Room() *Room {
	return p.room
}
