package list

import (
	"cmp"
	"fmt"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	kernel "github.com/Confidenceman02/scion-tools/pkg/list/internal"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"reflect"
)

type List[T any] interface {
	cons() *kernel.Cons[T, List[T]]
	Cmp(List[T]) int
	T() List[T]
}

func (c empty[T]) cons() *kernel.Cons[T, List[T]] {
	return &kernel.Cons[T, List[T]]{}
}
func (c *list[T]) cons() *kernel.Cons[T, List[T]] {
	return &kernel.Cons[T, List[T]]{}
}

func (c empty[T]) T() List[T] {
	return c
}
func (c *list[T]) T() List[T] {
	return c
}

type empty[T any] struct {
	*kernel.Cons[T, List[T]]
}

type list[T any] struct {
	*kernel.Cons[T, List[T]]
}

// Comparable

func (x *list[T]) Cmp(y List[T]) int {
	switch y := y.(type) {
	case empty[T]:
		return +1
	case *list[T]:
		// traverse conses until end of a list or a mismatch
		var ord = cmpHelp(x.Head, y.Head)
		var x1 *list[T] = x
		var y1 *list[T] = y
		for !IsEmpty(x1.Tail) && !IsEmpty(y1.Tail) && ord == 0 {
			switch x2 := x1.Tail.(type) {
			case *list[T]:
				switch y2 := y1.Tail.(type) {
				case *list[T]:
					x1 = x2
					y1 = y2
					ord = cmpHelp(x2.Head, y2.Head)
					continue
				}
			default:
				panic("unreachable")
			}
		}
		if IsEmpty(x1.Tail) && IsEmpty(y1.Tail) {
			return ord
		} else if !IsEmpty(x1.Tail) {
			return +1
		} else {
			return -1
		}
	default:
		var zero [0]T
		panic(
			fmt.Sprintf(
				"\nI was expecting a Comparable type, but instead got: \n    %v",
				reflect.TypeOf(zero).Elem(),
			),
		)
	}
}

func cmpHelp(x any, y any) int {
	switch x1 := x.(type) {
	case Int:
		switch y1 := y.(type) {
		case Int:
			return cmp.Compare(x1, y1)
		default:
			panic("Not an Int")
		}
	case Float:
		switch y1 := y.(type) {
		case Float:
			return cmp.Compare(x1, y1)
		default:
			panic("Not a Float")
		}
	case List[Int]:
		switch y1 := y.(type) {
		case List[Int]:
			return x1.Cmp(y1)
		default:
			panic("I was expecting a List[Int]")
		}
	case List[Float]:
		switch y1 := y.(type) {
		case List[Float]:
			return x1.Cmp(y1)
		default:
			panic("I was expecting a List[Float]")
		}
	default:
		panic(fmt.Sprintf("Cmp Not implemented for: %v", reflect.TypeOf(x1)))
	}
}

func (x empty[T]) Cmp(y List[T]) int {
	if reflect.DeepEqual(x, y) {
		return 0
	} else {
		return -1
	}
}

// CREATE

// Create a list with no elements
func Empty[T any]() List[T] {
	return empty[T]{}
}

// Create a list with only one element.
func Singleton[T any](val T) List[T] {
	return &list[T]{&kernel.Cons[T, List[T]]{Head: val, Tail: Empty[T]()}}
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
	return &list[T]{&kernel.Cons[T, List[T]]{Head: val, Tail: l}}
}

// TRANSFORM

func Map[A any, B any](f func(A) B, xs List[A]) List[B] {
	return Foldr(func(a A, b List[B]) List[B] { return Cons(f(a), b) }, Empty[B](), xs)
}

// Reduce a list from the left.
func Foldl[A any, B any](f func(A, B) B, acc B, ls List[A]) B {
	return ListWith(
		ls,
		func(List[A]) B { return acc },
		func(head A, tail List[A]) B { return Foldl(f, f(head, acc), tail) },
	)
}

// Reduce a list from the right.
func Foldr[A any, B any](fn func(A, B) B, acc B, ls List[A]) B {
	return foldrHelper(fn, acc, 0, ls)
}

func foldrHelper[A any, B any](fn func(A, B) B, acc B, ctr Int, ls List[A]) B {
	return ListWith(ls,
		func(List[A]) B { return acc },
		func(a A, r1 List[A]) B {
			return ListWith(
				r1,
				func(List[A]) B { return fn(a, acc) },
				func(b A, r2 List[A]) B {
					return ListWith(
						r2,
						func(List[A]) B { return fn(a, fn(b, acc)) },
						func(c A, r3 List[A]) B {
							return ListWith(
								r3,
								func(List[A]) B { return fn(a, fn(b, fn(c, acc))) },
								func(d A, r4 List[A]) B {
									var res B
									if Gt(ctr, 500) {
										res = Foldl(fn, acc, Reverse(r4))
									} else {
										res = foldrHelper(fn, acc, Add(ctr, 1), r4)
									}
									return fn(a, fn(b, fn(c, fn(d, res))))
								},
							)
						},
					)
				},
			)
		},
	)
}

// UTILITIES

// Determine the length of a list.
func Length[T any](ls List[T]) Int {
	return Foldl(func(_ T, y Int) Int { return y + 1 }, 0, ls)
}

// Reverse a list.
func Reverse[T any](ls List[T]) List[T] {
	return Foldl(Cons[T], Empty[T](), ls)
}

// Figure out whether a list contains a value.
func Member[T any](val T, l List[T]) bool {
	return Any(func(x T) bool { return Eq(x, val) }, l)
}

// Determine if all elements satisfy some test.
func All[T any](isOkay func(T) bool, l List[T]) bool {
	return Not(Any(ComposeL(Not, isOkay), l))
}

// Determine if any elements satisfy some test.
func Any[T any](isOkay func(T) bool, l List[T]) bool {
	return ListWith(
		l,
		func(List[T]) bool { return false },
		func(head T, tail List[T]) bool {
			if isOkay(head) {
				return true
			} else {
				return Any(isOkay, tail)
			}
		},
	)
}

// COMBINE

func Map2[A any, B any, result any](f func(A, B) result, xs List[A], ys List[B]) List[result] {
	return fromArray(map2Help(f, xs, ys, []result{}))
}

func map2Help[A any, B any, result any](f func(A, B) result, xs List[A], ys List[B], acc []result) []result {
	return ListWith(
		xs,
		func(List[A]) []result {
			return acc
		},
		func(x A, xt List[A]) []result {
			return ListWith(
				ys,
				func(List[B]) []result {
					return acc
				},
				func(y B, yt List[B]) []result {
					return map2Help(f, xt, yt, append(acc, f(x, y)))
				},
			)
		},
	)
}

// DECONSTRUCT

// Determine if a list is empty.
func IsEmpty[T any](l List[T]) bool {
	return ListWith(
		l,
		func(List[T]) bool { return true },
		func(_ T, _ List[T]) bool { return false },
	)
}

// Extract the first element of a list.
func Head[T any](l List[T]) maybe.Maybe[T] {
	return ListWith(
		l,
		func(List[T]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(head T, _ List[T]) maybe.Maybe[T] { return maybe.Just[T]{Value: head} },
	)
}

// Extract the rest of the list.
func Tail[T any](l List[T]) maybe.Maybe[List[T]] {
	return ListWith(
		l,
		func(List[T]) maybe.Maybe[List[T]] { return maybe.Nothing{} },
		func(_ T, tail List[T]) maybe.Maybe[List[T]] { return maybe.Just[List[T]]{Value: tail} },
	)
}

// PATTERN MATCH

func ListWith[T any, R any](l1 List[T], e func(List[T]) R, ht func(T, List[T]) R) R {
	switch l1 := l1.(type) {
	case empty[T]:
		return e(l1)
	case *list[T]:
		return ht(l1.Head, l1.Tail)
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

func fromArray[T any](arr []T) List[T] {
	var result List[T] = Empty[T]()
	for i := len(arr) - 1; i >= 0; i-- {
		result = Cons(arr[i], result)
	}
	return result
}
