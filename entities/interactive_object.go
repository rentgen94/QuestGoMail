package entities

type actionCodeGeneratorType func(args []string, items []*Item) (code int, err error)

type InteractiveObject interface {
	Namer
	Descriptable
	IsAccessible() bool
	SetAccessible(isAccessible bool)
	Interact(args []string, items []*Item) (string, error)
}

type boundInteractiveObject struct {
	name                string
	description         string
	isAccessible        bool
	room                *Room
	actionCodeGenerator actionCodeGeneratorType
}

func NewInteractiveObject(
	name string,
	description string,
	isAccessible bool,
	room *Room,
	actionCodeGenerator actionCodeGeneratorType,
) InteractiveObject {
	return &boundInteractiveObject{
		name:                name,
		description:         description,
		isAccessible:        isAccessible,
		room:                room,
		actionCodeGenerator: actionCodeGenerator,
	}
}

func (inter *boundInteractiveObject) Interact(args []string, items []*Item) (string, error) {
	var code, codeErr = inter.actionCodeGenerator(args, items)
	if codeErr != nil {
		return "", codeErr
	}

	return inter.room.PerformAction(code)
}

func (inter *boundInteractiveObject) IsAccessible() bool {
	return inter.isAccessible
}

func (inter *boundInteractiveObject) SetAccessible(isAccessible bool) {
	inter.isAccessible = isAccessible
}

func (inter *boundInteractiveObject) Name() string {
	return inter.name
}

func (inter *boundInteractiveObject) Description() string {
	return inter.description
}
