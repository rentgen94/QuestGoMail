package entities

type actType func(r *Room) (result InteractionResult, err error)

type Action struct {
	isAccessible bool
	name         string
	act          actType
}

func NewAction(name string, isAccessible bool, act actType) *Action {
	return &Action{
		name:         name,
		isAccessible: isAccessible,
		act:          act,
	}
}

func (a *Action) IsAccessible() bool {
	return a.isAccessible
}

func (a *Action) SetAccessible(isAccessible bool) {
	a.isAccessible = isAccessible
}

func (a *Action) Name() string {
	return a.name
}
