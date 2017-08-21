package entities

type Accessible interface {
	IsAccessible() bool
	SetAccessible(isAccessible bool)
}

type Namer interface {
	Name() string
}

type Descriptable interface {
	Description() string
}

type HasId interface {
	Id() int
}
