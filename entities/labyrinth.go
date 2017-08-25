package entities

type ActionsType map[int]*Action

type Labyrinth struct {
	id          int
	name        string
	description string
	rooms       RoomsType
	startRoom   *Room
	actions     ActionsType
	doors       DoorsType
}

func NewLabyrinth(id int, name string, description string) *Labyrinth {
	var lab = &Labyrinth{
		id:          id,
		name:        name,
		description: description,
		rooms:       make(RoomsType),
		actions:     make(ActionsType),
		doors:       make(DoorsType),
	}

	return lab
}

func (l *Labyrinth) Id() int {
	return l.id
}

func (l *Labyrinth) AddAction(action *Action) {
	l.actions[action.id] = action
	action.SetLabyrinth(l)
}

func (l *Labyrinth) SetActions(actions ActionsType) {
	for id, action := range actions {
		l.actions[id] = action
		action.SetLabyrinth(l)
	}
}

func (l *Labyrinth) Actions() ActionsType {
	return l.actions
}

func (l *Labyrinth) Rooms() RoomsType {
	return l.rooms
}

func (l *Labyrinth) SetRooms(rooms RoomsType) {
	l.rooms = rooms
}

func (l *Labyrinth) Doors() DoorsType {
	return l.doors
}

func (l *Labyrinth) SetDoors(doors DoorsType) {
	l.doors = doors
}

func (l *Labyrinth) StartRoom() *Room {
	return l.startRoom
}

func (l *Labyrinth) SetStartRoom(room *Room) {
	l.startRoom = room
}
