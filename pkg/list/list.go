// You can create a `List` from any Go slice with the `FromSlice` function. This module has a bunch of functions to help you work with them!
package list

import (
	"cmp"
	"fmt"
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/internal"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	. "github.com/Confidenceman02/scion-tools/pkg/tuple"
	"reflect"
	"slices"
)

type List[T any] interface {
	Cons() *internal.Cons_[T, List[T]]
	Cmp(basics.Comparable[List[T]]) int
	T() List[T]
}

func (c empty[T]) Cons() *internal.Cons_[T, List[T]] {
	return nil
}
func (c *list[T]) Cons() *internal.Cons_[T, List[T]] {
	return c.Cons_
}

func (c empty[T]) T() List[T] {
	return c
}
func (c *list[T]) T() List[T] {
	return c
}

type empty[T any] struct {
	*internal.Cons_[T, List[T]]
}

type list[T any] struct {
	*internal.Cons_[T, List[T]]
}

// Comparable

func (x *list[T]) Cmp(y basics.Comparable[List[T]]) int {
	switch y := y.T().(type) {
	case empty[T]:
		return +1
	case *list[T]:
		// traverse conses until end of a list or a mismatch
		var ord = cmpHelp(x.A, y.A)
		var x1 *list[T] = x
		var y1 *list[T] = y
		for !IsEmpty(x1.B) && !IsEmpty(y1.B) && ord == 0 {
			switch x2 := x1.B.(type) {
			case *list[T]:
				switch y2 := y1.B.(type) {
				case *list[T]:
					x1 = x2
					y1 = y2
					ord = cmpHelp(x2.A, y2.A)
					continue
				}
			default:
				panic("unreachable")
			}
		}
		if IsEmpty(x1.B) && IsEmpty(y1.B) {
			return ord
		} else if !IsEmpty(x1.B) {
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

func cmpHelp[T any](x T, y T) int {
	switch x1 := any(x).(type) {
	case List[T]:
		switch y1 := any(y).(type) {
		case List[T]:
			return x1.Cmp(y1)
		default:
			panic("I was expecting a List[T]")
		}
	default:
		return internal.CmpHelp(x, y)
	}
}

func (x empty[T]) Cmp(y basics.Comparable[List[T]]) int {
	if reflect.DeepEqual(x, y) {
		return 0
	} else {
		return -1
	}
}

// Create

// Create a list with no elements
func Empty[T any]() List[T] {
	return empty[T]{}
}

// Create a list with only one element.
func Singleton[T any](val T) List[T] {
	return &list[T]{&internal.Cons_[T, List[T]]{A: val, B: Empty[T]()}}
}

// Create a list with *n* copies of a value.
func Repeat[T any](n basics.Int, val T) List[T] {
	return repeatHelp(Empty[T](), n, val)

}

func repeatHelp[T any](result List[T], n basics.Int, val T) List[T] {
repeatHelpL:
	for {
		if n <= 0 {
			return result
		} else {
			tempResult, tempN, tempValue := Cons(val, result), n-1, val
			result = tempResult
			n = tempN
			val = tempValue
			continue repeatHelpL
		}
	}
}

// Create a list of numbers, every element increasing by one. You give the lowest and highest number that should be in the list.
func Range(low basics.Int, hi basics.Int) List[basics.Int] {
	return rangeHelp(low, hi, Empty[basics.Int]())
}

func rangeHelp(low basics.Int, hi basics.Int, ls List[basics.Int]) List[basics.Int] {
rangeHelpL:
	for {
		if cmp.Compare(low, hi) < 1 {
			tempLo, tempHi, tempLs := low, hi-1, Cons(hi, ls)
			low = tempLo
			hi = tempHi
			ls = tempLs
			continue rangeHelpL
		} else {
			return ls
		}
	}
}

// Add an element to the front of a list.
func Cons[T any](val T, l List[T]) List[T] {
	return &list[T]{&internal.Cons_[T, List[T]]{A: val, B: l}}
}

// Transform

// Apply a function to every element of a list.
func Map[A, B any](f func(A) B, xs List[A]) List[B] {
	return Foldr(func(a A, b List[B]) List[B] { return Cons(f(a), b) }, Empty[B](), xs)
}

// Same as map but the function is also applied to the index of each element (starting at zero).
func IndexedMap[A, B any](f func(basics.Int, A) B, xs List[A]) List[B] {
	return Map2(f, Range(0, basics.Sub(Length(xs), 1)), xs)
}

// Reduce a list from the left.
func Foldl[A, B any](f func(A, B) B, acc B, ls List[A]) B {
foldlL:
	for {
		if ls.Cons() == nil {
			return acc
		} else {
			var x = ls.Cons().A
			var xs = ls.Cons().B
			tempFunc, tempAcc, tempList := f, f(x, acc), xs
			f = tempFunc
			acc = tempAcc
			ls = tempList
			continue foldlL
		}
	}
}

// Reduce a list from the right.
func Foldr[A, B any](fn func(A, B) B, acc B, ls List[A]) B {
	return foldrHelper(fn, acc, 0, ls)
}

func foldrHelper[A, B any](fn func(A, B) B, acc B, ctr basics.Int, ls List[A]) B {
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
									if basics.Gt(ctr, basics.Int(500)) {
										res = Foldl(fn, acc, Reverse(r4))
									} else {
										res = foldrHelper(fn, acc, basics.Add(ctr, 1), r4)
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

// Keep elements that satisfy the test.
func Filter[T any](isGood func(T) bool, list List[T]) List[T] {
	return Foldr(func(x T, xs List[T]) List[T] {
		if isGood(x) {
			return Cons(x, xs)
		} else {
			return xs
		}
	}, Empty[T](), list)
}

// Filter out certain values. For example, maybe you have a bunch of strings from an
// untrusted source and you want to turn them into numbers:
func FilterMap[A, B any](f func(A) maybe.Maybe[B], xs List[A]) List[B] {
	return Foldr(func(a A, b List[B]) List[B] { return maybeCons(f, a, b) }, Empty[B](), xs)
}

func maybeCons[A, B any](f func(A) maybe.Maybe[B], mx A, xs List[B]) List[B] {
	return maybe.MaybeWith(
		f(mx),
		func(j maybe.Just[B]) List[B] { return Cons(j.Value, xs) },
		func(maybe.Nothing) List[B] { return xs },
	)
}

// Utilities

// Determine the length of a list.
func Length[T any](ls List[T]) basics.Int {
	return Foldl(func(_ T, y basics.Int) basics.Int { return y + 1 }, 0, ls)
}

// Reverse a list.
func Reverse[T any](ls List[T]) List[T] {
	return Foldl(Cons[T], Empty[T](), ls)
}

// Figure out whether a list contains a value.
func Member[T any](val T, l List[T]) bool {
	return Any(func(x T) bool { return basics.Eq(x, val) }, l)
}

// Determine if all elements satisfy some test.
func All[T any](isOkay func(T) bool, l List[T]) bool {
	return basics.Not(Any(basics.ComposeL(basics.Not, isOkay), l))
}

// Determine if any elements satisfy some test.
func Any[T any](isOkay func(T) bool, ls List[T]) bool {
anyL:
	for {
		if ls.Cons() == nil {
			return false
		} else {
			var x = ls.Cons().A
			var xs = ls.Cons().B
			if isOkay(x) {
				return true
			} else {
				tempIsOk, tempList := isOkay, xs
				isOkay = tempIsOk
				ls = tempList
				continue anyL
			}
		}
	}
}

// Find the maximum element in a non-empty list.
func Maximum[T basics.Comparable[T]](xs List[T]) maybe.Maybe[T] {
	return ListWith(
		xs,
		func(List[T]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(x T, xs List[T]) maybe.Maybe[T] {
			return maybe.Just[T]{Value: Foldl[T, T](basics.Max, x, xs).T()}
		},
	)
}

// Find the minimum element in a non-empty list.
func Minimum[T basics.Comparable[T]](xs List[T]) maybe.Maybe[T] {
	return ListWith(
		xs,
		func(List[T]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(x T, xs List[T]) maybe.Maybe[T] {
			return maybe.Just[T]{Value: Foldl[T, T](basics.Min, x, xs).T()}
		},
	)
}

// Get the sum of the list elements.
func Sum[T basics.Number](xs List[T]) T {
	return Foldl(basics.Add, 0, xs)
}

// Get the product of the list elements.
func Product[T basics.Number](xs List[T]) T {
	return Foldl(basics.Mul, 1, xs)
}

// Combine

// Put two lists together.
func Append[T any](xs List[T], ys List[T]) List[T] {
	return ListWith(xs,
		func(List[T]) List[T] { return ys },
		func(x T, xt List[T]) List[T] {
			return Foldr[T, List[T]](Cons, ys, xs)
		},
	)
}

// Concatenate a bunch of lists into a single list:
func Concat[T any](lists List[List[T]]) List[T] {
	return Foldr[List[T], List[T]](Append, Empty[T](), lists)
}

// Map a given function onto a list and flatten the resulting lists.
func ConcatMap[A, B any](f func(A) List[B], list List[A]) List[B] {
	return Concat(Map(f, list))
}

// Places the given value between all members of the given list.
func Intersperse[T any](sep T, xs List[T]) List[T] {
	return ListWith(
		xs,
		func(List[T]) List[T] { return Empty[T]() },
		func(hd T, t1 List[T]) List[T] {
			step := func(x T, rest List[T]) List[T] {
				return Cons(sep, Cons(x, rest))
			}
			spersed := Foldr(step, Empty[T](), t1)

			return Cons(hd, spersed)
		},
	)
}

// Combine two lists, combining them with the given function.
// If one list is longer, the extra elements are dropped.
func Map2[A any, B any, result any](f func(A, B) result, xs List[A], ys List[B]) List[result] {
	return FromSlice(map2Help(f, xs, ys))
}

func map2Help[A any, B any, result any](f func(A, B) result, xs List[A], ys List[B]) []result {
	var arr []result = []result{}
	for ; xs.Cons() != nil && ys.Cons() != nil; xs, ys = xs.Cons().B, ys.Cons().B {
		arr = append(arr, f(xs.Cons().A, ys.Cons().A))
	}
	return arr
}

func Map3[A, B, C, result any](f func(A, B, C) result, xs List[A], ys List[B], zs List[C]) List[result] {
	return FromSlice(map3Help(f, xs, ys, zs))
}

func map3Help[A any, B any, C any, result any](f func(A, B, C) result, xs List[A], ys List[B], zs List[C]) []result {
	var arr []result = []result{}
	for ; xs.Cons() != nil && ys.Cons() != nil && zs.Cons() != nil; xs, ys, zs = xs.Cons().B, ys.Cons().B, zs.Cons().B {
		arr = append(arr, f(xs.Cons().A, ys.Cons().A, zs.Cons().A))
	}
	return arr
}

func Map4[A, B, C, D, result any](f func(A, B, C, D) result, xs List[A], ys List[B], zs List[C], ws List[D]) List[result] {
	return FromSlice(map4Help(f, xs, ys, zs, ws))
}

func map4Help[A, B, C, D, result any](f func(A, B, C, D) result, ws List[A], xs List[B], ys List[C], zs List[D]) []result {
	var arr []result = []result{}
	for ; ws.Cons() != nil && xs.Cons() != nil && ys.Cons() != nil && zs.Cons() != nil; ws, xs, ys, zs = ws.Cons().B, xs.Cons().B, ys.Cons().B, zs.Cons().B {
		arr = append(arr, f(ws.Cons().A, xs.Cons().A, ys.Cons().A, zs.Cons().A))
	}
	return arr
}

func Map5[A, B, C, D, E, result any](f func(A, B, C, D, E) result, vs List[A], ws List[B], xs List[C], ys List[D], zs List[E]) List[result] {
	return FromSlice(map5Help(f, vs, ws, xs, ys, zs))
}

func map5Help[A, B, C, D, E, result any](f func(A, B, C, D, E) result, vs List[A], ws List[B], xs List[C], ys List[D], zs List[E]) []result {
	var arr []result = []result{}
	for ; vs.Cons() != nil && ws.Cons() != nil && xs.Cons() != nil && ys.Cons() != nil && zs.Cons() != nil; vs, ws, xs, ys, zs = vs.Cons().B, ws.Cons().B, xs.Cons().B, ys.Cons().B, zs.Cons().B {
		arr = append(arr, f(vs.Cons().A, ws.Cons().A, xs.Cons().A, ys.Cons().A, zs.Cons().A))
	}
	return arr
}

// Sort

// Sort values from lowest to highest.
func Sort[T basics.Comparable[T]](xs List[T]) List[T] {
	slc := ToSlice(xs)
	slices.SortFunc(
		slc,
		func(a, b T) int {
			return a.Cmp(b)
		},
	)
	return FromSlice(slc)
}

// Sort values by a derived property.
func SortBy[A basics.Comparable[A]](f func(A) A, xs List[A]) List[A] {
	slc := ToSlice(xs)
	slices.SortFunc(
		slc,
		func(a, b A) int {
			return f(a).Cmp(f(b))
		},
	)
	return FromSlice(slc)
}

// Sort values with a custom comparison function.
func SortWith[A any](f func(a A, b A) basics.Order, xs List[A]) List[A] {
	slc := ToSlice(xs)
	slices.SortFunc(slc, func(a, b A) int {
		ord := f(a, b)
		switch ord.(type) {
		case basics.EQ:
			return 0
		case basics.LT:
			return -1
		default:
			return 1
		}
	})
	return FromSlice(slc)
}

// Deconstruct

// Determine if a list is empty.
func IsEmpty[T any](l List[T]) bool {
	return l.Cons() == nil
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

// Take the first n members of a list.
func Take[T any](n basics.Int, list List[T]) List[T] {
	return takeFast(0, n, list)
}

func takeFast[A any](ctr basics.Int, n basics.Int, list List[A]) List[A] {
	// This looks gross because it's an analogue to the compiled Elm kernel code.
	// There is almost definitely a cleaner way to do this but it is
	// here mostly for examining how kernel code translates to Go.
	if n <= 0 {
		return Empty[A]()
	} else {
		cns := list.Cons()
	loop1:
		for {
		loop2:
			for {
				if cns == nil {
					return list
				} else if cns.B.Cons() != nil {
					switch n {
					case 1:
						break loop1
					case 2:
						x := cns.A
						val2 := cns.B.Cons()
						y := val2.A
						return FromSlice([]A{x, y})
					case 3:
						if cns.B.Cons() != nil && cns.B.Cons().B.Cons() != nil {
							x := cns.A
							val2 := cns.B.Cons()
							y := val2.A
							val3 := val2.B.Cons()
							z := val3.A
							return FromSlice([]A{x, y, z})
						} else {
							break loop2
						}
					default:
						if cns.B.Cons() != nil && cns.B.Cons().B.Cons() != nil && cns.B.Cons().B.Cons().B.Cons() != nil {
							x := cns.A
							y := cns.B.Cons().A
							z := cns.B.Cons().B.Cons().A
							w := cns.B.Cons().B.Cons().B.Cons().A
							tl := cns.B.Cons().B.Cons().B.Cons().B
							if ctr > 1000 {
								return Cons(x, Cons(y, Cons(z, Cons(w, takeTailRec(n-4, tl)))))
							} else {
								return Cons(x, Cons(y, Cons(z, Cons(w, takeFast(ctr+1, n-4, tl)))))
							}
						}
					}
				} else {
					if n == 1 {
						break loop1
					} else {
						break loop2
					}
				}
			}
			return list
		}
		return FromSlice([]A{cns.A})
	}
}

func takeTailRec[A any](n basics.Int, list List[A]) List[A] {
	return Reverse(takeReverse(n, list, Empty[A]()))
}

func takeReverse[A any](n basics.Int, list List[A], kept List[A]) List[A] {
takeReverseL:
	for {
		if n <= 0 {
			return kept
		} else {
			if list.Cons() == nil {
				return kept
			} else {
				var x = list.Cons().A
				var xs = list.Cons().B
				tempN, tempList, tempKept := n-1, xs, Cons(x, kept)
				n = tempN
				list = tempList
				kept = tempKept
				continue takeReverseL
			}
		}
	}
}

// Drop the first n members of a list.
func Drop[T any](n basics.Int, list List[T]) List[T] {
	if n <= 0 {
		return list
	} else {
		return ListWith(
			list,
			func(List[T]) List[T] { return list },
			func(x T, xs List[T]) List[T] {
				return Drop(n-1, xs)
			},
		)
	}
}

// Partition a list based on some test. The first list contains all values
// that satisfy the test, and the second list contains all the value that do not.
func Partition[A any](pred func(A) bool, list List[A]) Tuple2[List[A], List[A]] {
	step := func(x A, tf Tuple2[List[A], List[A]]) Tuple2[List[A], List[A]] {
		if pred(x) {
			return Pair(Cons(x, First(tf)), Second(tf))
		} else {
			return Pair(First(tf), Cons(x, Second(tf)))
		}
	}
	return Foldr(step, Pair(Empty[A](), Empty[A]()), list)
}

// Decompose a list of tuples into a tuple of lists.
func Unzip[A, B any](pairs List[Tuple2[A, B]]) Tuple2[List[A], List[B]] {
	step := func(tuple Tuple2[A, B], acc Tuple2[List[A], List[B]]) Tuple2[List[A], List[B]] {
		return Pair(Cons(First(tuple), First(acc)), Cons(Second(tuple), Second(acc)))
	}
	return Foldr(step, Pair(Empty[A](), Empty[B]()), pairs)
}

// Pattern Match

func ListWith[T any, R any](l1 List[T], e func(List[T]) R, ab func(T, List[T]) R) R {
	if l1.Cons() == nil {
		return e(l1)
	} else {
		return ab(l1.Cons().A, l1.Cons().B)
	}
}

// Utils

// Create a List from a Go slice
func FromSlice[T any](arr []T) List[T] {
	var result List[T] = Empty[T]()
	for i := len(arr) - 1; i >= 0; i-- {
		result = Cons(arr[i], result)
	}
	return result
}
func FromSliceMap[A any, B any](f func(A) B, arr []A) List[B] {
	var result List[B] = Empty[B]()
	for i := len(arr) - 1; i >= 0; i-- {
		result = Cons(f(arr[i]), result)
	}
	return result
}

func ToSlice[T any](xs List[T]) []T {
	var arr []T = []T{}
	for ; xs.Cons() != nil; xs = xs.Cons().B {
		arr = append(arr, xs.Cons().A)
	}
	return arr
}

func ToSliceMap[A any, B any](f func(A) B, xs List[A]) []B {
	var arr []B = []B{}
	for ; xs.Cons() != nil; xs = xs.Cons().B {
		arr = append(arr, f(xs.Cons().A))
	}
	return arr
}
