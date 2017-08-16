package entities

type InteractiveObject interface {
	IsAccessible() bool
	SetAccessible(isAccessible bool)
	Interact(args []string, items []*Item) error
}
