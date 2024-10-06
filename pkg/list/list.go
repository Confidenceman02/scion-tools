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

// Create a list with only one element.
func Singleton[T any](val T) List[T] {
	return &list[T]{&cons[T]{val, nil}}
}

// Add an element to the front of a list.
func Cons[T any](val T, l List[T]) List[T] {
	return &list[T]{&cons[T]{val, l._list()}}
}

// Create a list with *n* copies of a value.
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

// Create a list of numbers, every element increasing by one. You give the lowest and highest number that should be in the list.
func Range(low basics.Int, hi basics.Int) List[basics.Int] {
	return rangeHelp(low, hi, Empty[basics.Int]())
}

func rangeHelp(low basics.Int, hi basics.Int, ls List[basics.Int]) List[basics.Int] {
	if low <= hi {
		return rangeHelp(low, hi-1, Cons(hi, ls))
	} else {
		return ls
	}
}
