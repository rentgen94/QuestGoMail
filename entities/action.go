package entities

type roomsType map[int]*Room
type actType func(labyrinth *Labyrinth) (result InteractionResult, err error)

type Action struct {
	id int
	isAccessible bool
	name         string
	labyrinth *Labyrinth
	act          actType
}

func NewAction(name string, labyrinth *Labyrinth, isAccessible bool, act actType) *Action {
	return &Action{
		name:         name,
		isAccessible: isAccessible,
		act:          act,
		labyrinth: labyrinth,
	}
}


// TODO maybe Action does not need isAccessible member
func (a *Action) IsAccessible() bool {
	return a.isAccessible
}

func (a *Action) SetAccessible(isAccessible bool) {
	a.isAccessible = isAccessible
}

func (a *Action) Name() string {
	return a.name
}

func (a *Action) Act() (result InteractionResult, err error) {
	return a.act(a.labyrinth)
}
