package entities

type ActionsType map[int]*Action

type Labyrinth struct {
	id        int
	rooms     RoomsType
	startRoom *Room
	actions   ActionsType
}

func NewLabyrinth(rooms RoomsType, startRoom *Room, actions ActionsType) *Labyrinth {
	var lab = &Labyrinth{
		rooms:     rooms,
		startRoom: startRoom,
		actions:   actions,
	}

	for key := range lab.actions {
		lab.actions[key].SetLabyrinth(lab)
	}

	return lab
}

func (l *Labyrinth) Rooms() RoomsType {
	return l.rooms
}

func (l *Labyrinth) StartRoom() *Room {
	return l.startRoom
}
