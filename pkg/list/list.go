package list

import (
	"fmt"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"reflect"
)

type List[T any] interface {
	_consList() consList
}

type consList struct{}

func Empty[T any]() List[T] {
	return empty[T]{}
}

func (cl consList) _consList() consList {
	return cl
}

type empty[T any] struct {
	consList
}

type list[T any] struct {
	consList
	_cons *cons[T]
}

type cons[T any] struct {
	head T
	tail List[T]
}

// Create a list with only one element.
func Singleton[T any](val T) List[T] {
	return &list[T]{consList{}, &cons[T]{val, Empty[T]()}}
}

// Add an element to the front of a list.
func Cons[T any](val T, l List[T]) List[T] {
	return &list[T]{consList{}, &cons[T]{val, l}}
}

// Create a list with *n* copies of a value.
func Repeat[T any](n Int, val T) List[T] {
	return repeatHelp(Empty[T](), n, val)

}

func repeatHelp[T any](result List[T], n Int, val T) List[T] {
	if n <= 0 {
		return result
	} else {
		return repeatHelp(Cons(val, result), n-1, val)
	}
}

// Create a list of numbers, every element increasing by one. You give the lowest and highest number that should be in the list.
func Range(low Int, hi Int) List[Int] {
	return rangeHelp(low, hi, Empty[Int]())
}

func rangeHelp(low Int, hi Int, ls List[Int]) List[Int] {
	if low <= hi {
		return rangeHelp(low, hi-1, Cons(hi, ls))
	} else {
		return ls
	}
}

// DECONSTRUCT

// Determine if a list is empty.
func IsEmpty[T any](l List[T]) bool {
	return ListWith(
		l,
		func(empty[T]) bool { return true },
		func(*list[T]) bool { return false },
	)
}

// Extract the first element of a list.
func Head[T any](l List[T]) maybe.Maybe[T] {
	return ListWith(
		l,
		func(empty[T]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(l *list[T]) maybe.Maybe[T] { return maybe.Just[T]{Value: l._cons.head} },
	)
}

// Extract the rest of the list.
func Tail[T any](l List[T]) maybe.Maybe[List[T]] {
	return ListWith(
		l,
		func(empty[T]) maybe.Maybe[List[T]] { return maybe.Nothing{} },
		func(l *list[T]) maybe.Maybe[List[T]] { return maybe.Just[List[T]]{Value: l._cons.tail} },
	)
}

// PATTERN MATCH

func ListWith[T any, R any](l1 List[T], e func(empty[T]) R, ne func(*list[T]) R) R {
	switch l1 := l1.(type) {
	case empty[T]:
		return e(l1)
	case *list[T]:
		return ne(l1)
	default:
		var zero [0]T
		panic(
			fmt.Sprintf(
				"\nI was expecting a type of: \n    *list.list[%v]\n\nBut instead got a\n    %v\n",
				reflect.TypeOf(zero).Elem(),
				reflect.TypeOf(l1),
			),
		)
	}
}
