package internal

type Cons[A any, B any] struct {
	Head A
	Tail B
}

func List_Cons[A any, B any](head A, tail B) *Cons[A, B] {
	return &Cons[A, B]{head, tail}
}
