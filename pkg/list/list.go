package list

import (
	"fmt"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/internal"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	. "github.com/Confidenceman02/scion-tools/pkg/tuple"
	"reflect"
	"slices"
)

type _cons[A any] struct {
	a A
	b List[A]
}

type List[T any] interface {
	cons() *_cons[T]
	Cmp(Comparable[List[T]]) int
	T() List[T]
}

func (c empty[T]) cons() *_cons[T] {
	return nil
}
func (c *list[T]) cons() *_cons[T] {
	return c._cons
}

func (c empty[T]) T() List[T] {
	return c
}
func (c *list[T]) T() List[T] {
	return c
}

type empty[T any] struct {
	*_cons[T]
}

type list[T any] struct {
	*_cons[T]
}

// Comparable

func (x *list[T]) Cmp(y Comparable[List[T]]) int {
	switch y := y.T().(type) {
	case empty[T]:
		return +1
	case *list[T]:
		// traverse conses until end of a list or a mismatch
		var ord = cmpHelp(x.a, y.a)
		var x1 *list[T] = x
		var y1 *list[T] = y
		for !IsEmpty(x1.b) && !IsEmpty(y1.b) && ord == 0 {
			switch x2 := x1.b.(type) {
			case *list[T]:
				switch y2 := y1.b.(type) {
				case *list[T]:
					x1 = x2
					y1 = y2
					ord = cmpHelp(x2.a, y2.a)
					continue
				}
			default:
				panic("unreachable")
			}
		}
		if IsEmpty(x1.b) && IsEmpty(y1.b) {
			return ord
		} else if !IsEmpty(x1.b) {
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

func (x empty[T]) Cmp(y Comparable[List[T]]) int {
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
	return &list[T]{&_cons[T]{a: val, b: Empty[T]()}}
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
	return &list[T]{&_cons[T]{a: val, b: l}}
}

// Transform

func Map[A any, B any](f func(A) B, xs List[A]) List[B] {
	return Foldr(func(a A, b List[B]) List[B] { return Cons(f(a), b) }, Empty[B](), xs)
}

func IndexedMap[A any, B any](f func(Int, A) B, xs List[A]) List[B] {
	return Map2(f, Range(0, Sub(Length(xs), 1)), xs)
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
									if Gt(ctr, Int(500)) {
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

func FilterMap[A any, B any](f func(A) maybe.Maybe[B], xs List[A]) List[B] {
	return Foldr(func(a A, b List[B]) List[B] { return maybeCons(f, a, b) }, Empty[B](), xs)
}

func maybeCons[A any, B any](f func(A) maybe.Maybe[B], mx A, xs List[B]) List[B] {
	return maybe.MaybeWith(
		f(mx),
		func(j maybe.Just[B]) List[B] { return Cons(j.Value, xs) },
		func(maybe.Nothing) List[B] { return xs },
	)
}

// Utilities

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

// Find the maximum element in a non-empty list.
func Maximum[T any](xs List[Comparable[T]]) maybe.Maybe[T] {
	return ListWith(
		xs,
		func(List[Comparable[T]]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(x Comparable[T], xs List[Comparable[T]]) maybe.Maybe[T] {
			return maybe.Just[T]{Value: Foldl[Comparable[T], Comparable[T]](Max, x, xs).T()}
		},
	)
}

// Find the maximum element in a non-empty list of generic elements.
// If elements do not implement the Comparable[T] interface the function will panic.
func Maximum_UNSAFE[T any](xs List[T]) maybe.Maybe[T] {
	return ListWith(
		xs,
		func(List[T]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(x T, xt List[T]) maybe.Maybe[T] {
			restSlice := toSlice(xt)

			if comp1, ok := any(x).(Comparable[T]); ok {
				var ret Comparable[T] = comp1
				for idx := 0; idx < len(restSlice); idx++ {
					if comp2, ok := any(restSlice[idx]).(Comparable[T]); ok {
						ret = Max(ret, comp2)
					} else {
						panic("Cannot find Maximum_UNSAFE on non-comparable types")
					}
				}
				return maybe.Just[T]{Value: ret.T()}
			} else {
				panic("Cannot find Maximum_UNSAFE on non-comparable types")
			}
		},
	)
}

// Find the minimum element in a non-empty list.
func Minimum[T any](xs List[Comparable[T]]) maybe.Maybe[T] {
	return ListWith(
		xs,
		func(List[Comparable[T]]) maybe.Maybe[T] { return maybe.Nothing{} },
		func(x Comparable[T], xs List[Comparable[T]]) maybe.Maybe[T] {
			return maybe.Just[T]{Value: Foldl[Comparable[T], Comparable[T]](Min, x, xs).T()}
		},
	)
}

// Get the sum of the list elements.
func Sum[T Number](xs List[T]) T {
	return Foldl(Add, 0, xs)
}

// Get the product of the list elements.
func Product[T Number](xs List[T]) T {
	return Foldl(Mul, 1, xs)
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
func ConcatMap[A any, B any](f func(A) List[B], list List[A]) List[B] {
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
	return fromSlice(map2Help(f, xs, ys, []result{}))
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

func Map3[A any, B any, C any, result any](f func(A, B, C) result, xs List[A], ys List[B], zs List[C]) List[result] {
	return fromSlice(map3Help(f, xs, ys, zs, []result{}))
}

func map3Help[A any, B any, C any, result any](f func(A, B, C) result, xs List[A], ys List[B], zs List[C], acc []result) []result {
	return ListWith(
		xs,
		func(List[A]) []result { return acc },
		func(x A, xs List[A]) []result {
			return ListWith(
				ys,
				func(List[B]) []result { return acc },
				func(y B, ys List[B]) []result {
					return ListWith(
						zs,
						func(List[C]) []result { return acc },
						func(z C, zs List[C]) []result {
							return map3Help(f, xs, ys, zs, append(acc, f(x, y, z)))
						},
					)
				},
			)
		},
	)
}

func Map4[A any, B any, C any, D any, result any](f func(A, B, C, D) result, xs List[A], ys List[B], zs List[C], ws List[D]) List[result] {
	return fromSlice(map4Help(f, xs, ys, zs, ws, []result{}))
}

func map4Help[A any, B any, C any, D any, result any](f func(A, B, C, D) result, ws List[A], xs List[B], ys List[C], zs List[D], acc []result) []result {
	return ListWith(
		ws,
		func(List[A]) []result { return acc },
		func(a A, as List[A]) []result {
			return ListWith(
				xs,
				func(List[B]) []result { return acc },
				func(b B, bs List[B]) []result {
					return ListWith(
						ys,
						func(List[C]) []result { return acc },
						func(c C, cs List[C]) []result {
							return ListWith(
								zs,
								func(List[D]) []result { return acc },
								func(d D, ds List[D]) []result {
									return map4Help(f, as, bs, cs, ds, append(acc, f(a, b, c, d)))
								},
							)
						},
					)
				},
			)
		},
	)
}

func Map5[A any, B any, C any, D any, E any, result any](f func(A, B, C, D, E) result, vs List[A], ws List[B], xs List[C], ys List[D], zs List[E]) List[result] {
	return fromSlice(map5Help(f, vs, ws, xs, ys, zs, []result{}))
}

func map5Help[A any, B any, C any, D any, E any, result any](f func(A, B, C, D, E) result, vs List[A], ws List[B], xs List[C], ys List[D], zs List[E], acc []result) []result {
	return ListWith(
		vs,
		func(List[A]) []result { return acc },
		func(a A, as List[A]) []result {
			return ListWith(
				ws,
				func(List[B]) []result { return acc },
				func(b B, bs List[B]) []result {
					return ListWith(
						xs,
						func(List[C]) []result { return acc },
						func(c C, cs List[C]) []result {
							return ListWith(
								ys,
								func(List[D]) []result { return acc },
								func(d D, ds List[D]) []result {
									return ListWith(
										zs,
										func(List[E]) []result { return acc },
										func(e E, es List[E]) []result {
											return map5Help(f, as, bs, cs, ds, es, append(acc, f(a, b, c, d, e)))
										},
									)
								},
							)
						},
					)
				},
			)
		},
	)
}

// Sort

// Sort values from lowest to highest.
func Sort[T any](xs List[Comparable[T]]) List[Comparable[T]] {
	return SortBy_UNSAFE(Identity, xs)
}

// Sort generic values from lowest to highest.
// This function will panic if T is not a Comparable[T]
func Sort_UNSAFE[T any](xs List[T]) List[T] {
	return SortBy_UNSAFE(Identity, xs)
}

// Sort values by a derived property.
func SortBy[A any](f func(A) Comparable[A], xs List[A]) List[A] {
	slc := toSlice(xs)
	slices.SortFunc(
		slc,
		func(a, b A) int {
			return f(a).Cmp(f(b))
		},
	)
	return fromSlice(slc)
}

// Sort values by a derived property.
// This function will panic if the passed in func doesn't return a Comparable.
func SortBy_UNSAFE[A any](f func(A) A, xs List[A]) List[A] {
	slc := toSlice(xs)
	slices.SortFunc(
		slc,
		func(a, b A) int {
			if comp1, ok := any(f(a)).(Comparable[A]); ok {
				if comp2, ok := any(f(b)).(Comparable[A]); ok {
					return comp1.Cmp(comp2)
				} else {
					panic("I was expecting a Comparable type")
				}
			} else {
				panic("I was expecting a Comparable type")
			}
		},
	)
	return fromSlice(slc)
}

func SortWith[A any](f func(a A, b A) Order, xs List[A]) List[A] {
	slc := toSlice(xs)
	slices.SortFunc(slc, func(a, b A) int {
		ord := f(a, b)
		switch ord.(type) {
		case EQ:
			return 0
		case LT:
			return -1
		default:
			return 1
		}
	})
	return fromSlice(slc)
}

// Deconstruct

// Determine if a list is empty.
func IsEmpty[T any](l List[T]) bool {
	return l.cons() == nil
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
func Take[T any](n Int, list List[T]) List[T] {
	return takeFast(0, n, list)
}

func takeFast[A any](ctr Int, n Int, list List[A]) List[A] {
	// This looks gross because it's an analogue to the compiled Elm kernel code.
	// There is almost definitely a cleaner way to do this but it is
	// here mostly for examining how kernel code translates to Go.
	if n <= 0 {
		return Empty[A]()
	} else {
		cns := list.cons()
	loop1:
		for {
		loop2:
			for {
				if cns == nil {
					return list
				} else if cns.b.cons() != nil {
					switch n {
					case 1:
						break loop1
					case 2:
						x := cns.a
						val2 := cns.b.cons()
						y := val2.a
						return fromSlice([]A{x, y})
					case 3:
						if cns.b.cons() != nil && cns.b.cons().b.cons() != nil {
							x := cns.a
							val2 := cns.b.cons()
							y := val2.a
							val3 := val2.b.cons()
							z := val3.a
							return fromSlice([]A{x, y, z})
						} else {
							break loop2
						}
					default:
						if cns.b.cons() != nil && cns.b.cons().b.cons() != nil && cns.b.cons().b.cons().b.cons() != nil {
							x := cns.a
							y := cns.b.cons().a
							z := cns.b.cons().b.cons().a
							w := cns.b.cons().b.cons().b.cons().a
							tl := cns.b.cons().b.cons().b.cons().b
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
		return fromSlice([]A{cns.a})
	}
}

func takeTailRec[A any](n Int, list List[A]) List[A] {
	return Reverse(takeReverse(n, list, Empty[A]()))
}

func takeReverse[A any](n Int, list List[A], kept List[A]) List[A] {
	if n <= 0 {
		return kept
	} else {
		return ListWith(
			list,
			func(List[A]) List[A] { return kept },
			func(x A, xs List[A]) List[A] {
				return takeReverse(n-1, xs, Cons(x, kept))
			},
		)
	}
}

func Drop[T any](n Int, list List[T]) List[T] {
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
	if l1.cons() == nil {
		return e(l1)
	} else {
		return ab(l1.cons().a, l1.cons().b)
	}
}

// Utils

func fromSlice[T any](arr []T) List[T] {
	var result List[T] = Empty[T]()
	for i := len(arr) - 1; i >= 0; i-- {
		result = Cons(arr[i], result)
	}
	return result
}

func toSlice[T any](xs List[T]) []T {
	return toSliceHelp(xs, []T{})
}

func toSliceHelp[T any](xs List[T], arr []T) []T {
	return ListWith(
		xs,
		func(List[T]) []T { return arr },
		func(x T, xr List[T]) []T {
			return toSliceHelp(xr, append(arr, x))
		})
}
