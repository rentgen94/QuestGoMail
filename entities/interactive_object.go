package entities

import "errors"

const (
	objectNotAccessible = "Object not accessible"
	actionNotAssigned   = "Action not assigned"
)

type InteractiveObject interface {
	Namer
	Descriptable
	HasId
	IsAccessible() bool
	SetAccessible(isAccessible bool)
	Interact(args []string, items []Item) (result InteractionResult, err error)
	SetAction(action *Action)
}

type InputChecker func(args []string, items []Item) error

type boundInteractiveObject struct {
	id           int
	name         string
	description  string
	isAccessible bool
	room         *Room
	checker      InputChecker
	action       *Action
}

func NewInteractiveObject(
	id int,
	name string,
	description string,
	isAccessible bool,
	checker InputChecker,
	action *Action,
) InteractiveObject {
	return &boundInteractiveObject{
		name:         name,
		id:           id,
		description:  description,
		isAccessible: isAccessible,
		checker:      checker,
		action:       action,
	}
}

func (inter *boundInteractiveObject) Interact(args []string, items []Item) (InteractionResult, error) {
	if !inter.isAccessible {
		return ContinueResult(""), errors.New(objectNotAccessible)
	}

	var checkErr = inter.checker(args, items)
	if checkErr != nil {
		return ContinueResult(""), checkErr
	}

	if inter.action == nil {
		return ContinueResult(""), errors.New(actionNotAssigned)
	}

	return inter.action.Act()
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

func (inter *boundInteractiveObject) SetAction(action *Action) {
	inter.action = action
}
