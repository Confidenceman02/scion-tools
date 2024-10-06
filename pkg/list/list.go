package list

import "github.com/Confidenceman02/scion-tools/pkg/basics"

type List[T any] interface {
	_list() *list[T]
}

func Empty[T any]() List[T] {
	return &list[T]{}
}

func (l *list[T]) _list() *list[T] {
	return l
}

type list[T any] struct {
	_cons *cons[T]
}

type cons[T any] struct {
	head T
	tail *list[T]
}

func Singleton[T any](val T) List[T] {
	return &list[T]{&cons[T]{val, nil}}
}

func Cons[T any](val T, l List[T]) List[T] {
	return &list[T]{&cons[T]{val, l._list()}}
}

func Repeat[T any](n basics.Int, val T) List[T] {
	return repeatHelp(Empty[T](), n, val)

}

func repeatHelp[T any](result List[T], n basics.Int, val T) List[T] {
	if n <= 0 {
		return result
	} else {
		return repeatHelp(Cons(val, result), n-1, val)
	}
}
