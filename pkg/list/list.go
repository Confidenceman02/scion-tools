package list

import (
	"fmt"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	. "github.com/Confidenceman02/scion-tools/pkg/maybe"
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

// Add an element to the front of a list.
func Cons[T any](val T, l List[T]) List[T] {
	return &list[T]{consList{}, &cons[T]{val, l}}
}

// TRANSFORM

// Reduce a list from the left.
func Foldl[A any, B any](f func(A, B) B, acc B, ls List[A]) B {
	return ListWith(
		ls,
		func(List[A]) B { return acc },
		func(head A, tail List[A]) B { return Foldl(f, f(head, acc), tail) },
	)
}

// UTILITY

// Determine the length of a list.
func Length(ls List[any]) Int {
	return Foldl(func(_, y Int) Int { return y + 1 }, 0, ls)
}

// Reverse a list.
func Reverse[T any](ls List[T]) List[T] {
	return Foldl(Cons[T], Empty[T](), ls)
}

// DECONSTRUCT

// Determine if a list is empty.
func IsEmpty[T any](l List[T]) bool {
	return ListWith(
		l,
		func(List[T]) bool { return true },
		func(head T, tail List[T]) bool { return false },
	)
}

// Extract the first element of a list.
func Head[T any](l List[T]) Maybe[T] {
	return ListWith(
		l,
		func(List[T]) Maybe[T] { return Nothing{} },
		func(head T, tail List[T]) Maybe[T] { return Just[T]{Value: head} },
	)
}

// Extract the rest of the list.
func Tail[T any](l List[T]) Maybe[List[T]] {
	return ListWith(
		l,
		func(List[T]) Maybe[List[T]] { return Nothing{} },
		func(head T, tail List[T]) Maybe[List[T]] { return Just[List[T]]{Value: tail} },
	)
}

// PATTERN MATCH

func ListWith[T any, R any](l1 List[T], e func(List[T]) R, ht func(T, List[T]) R) R {
	switch l1 := l1.(type) {
	case empty[T]:
		return e(l1)
	case *list[T]:
		return ht(l1._cons.head, l1._cons.tail)
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
