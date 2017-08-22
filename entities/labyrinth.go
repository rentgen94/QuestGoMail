package entities

type actionsType map[int]*Action

type Labyrinth struct {
	id int
	rooms roomsType
	startRoom *Room
	actions actionsType
}
