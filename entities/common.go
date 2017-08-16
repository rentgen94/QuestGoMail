package entities

type Accessible interface {
	IsAccessible() bool
	SetAccessible(isAccessible bool)
}
