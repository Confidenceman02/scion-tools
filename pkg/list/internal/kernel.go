package internal

type Cons[A any, B any] struct {
	Head A
	Tail B
}
