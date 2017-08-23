package entities

type RoomsType map[int]*Room
type ActType func(labyrinth *Labyrinth) (result InteractionResult, err error)

type Action struct {
	id        int
	name      string
	labyrinth *Labyrinth
	act       ActType
}

func NewAction(name string, act ActType) *Action {
	return &Action{
		name:      name,
		act:       act,
		labyrinth: nil,
	}
}

func (a *Action) Name() string {
	return a.name
}

func (a *Action) Act() (result InteractionResult, err error) {
	return a.act(a.labyrinth)
}

func (a *Action) GetLabyrinth() *Labyrinth {
	return a.labyrinth
}

func (a *Action) SetLabyrinth(labyrinth *Labyrinth) {
	a.labyrinth = labyrinth
}
