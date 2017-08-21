package entities

import "errors"

const (
	objectNotAccessible = "Object not accessible"
)

type actionCodeGeneratorType func(args []string, items []Item) (code int, err error)

type InteractiveObject interface {
	Namer
	Descriptable
	HasId
	IsAccessible() bool
	SetAccessible(isAccessible bool)
	Interact(args []string, items []Item) (result InteractionResult, err error)
}

type boundInteractiveObject struct {
	id                  int
	name                string
	description         string
	isAccessible        bool
	room                *Room
	actionCodeGenerator actionCodeGeneratorType
}

func NewInteractiveObject(
	name string,
	id int,
	description string,
	isAccessible bool,
	room *Room,
	actionCodeGenerator actionCodeGeneratorType,
) InteractiveObject {
	return &boundInteractiveObject{
		name:                name,
		id:                  id,
		description:         description,
		isAccessible:        isAccessible,
		room:                room,
		actionCodeGenerator: actionCodeGenerator,
	}
}

func (inter *boundInteractiveObject) Interact(args []string, items []Item) (InteractionResult, error) {
	if !inter.isAccessible {
		return ContinueResult(""), errors.New(objectNotAccessible)
	}

	var code, codeErr = inter.actionCodeGenerator(args, items)
	if codeErr != nil {
		return ContinueResult(""), codeErr
	}

	return inter.room.PerformAction(code)
}

func (inter *boundInteractiveObject) IsAccessible() bool {
	return inter.isAccessible
}

func (inter *boundInteractiveObject) Id() int {
	return inter.id
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
