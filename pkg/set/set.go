//	A set of unique values. The values can be any comparable type. This includes Int, Float, Char, String, and tuples or lists of comparable types.
//
// Insert, remove, and query operations all take O(log n) time.
package set

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/dict"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
)

type Set[K Comparable[K]] interface {
	set_() *set[K]
}

/*
Retrieve the internal set
*/
func (d *set[K]) set_() *set[K] {
	return d
}

type set[T Comparable[T]] struct {
	d dict.Dict[T, struct{}]
}

// Create an empty set.
func Empty[K Comparable[K]]() Set[K] {
	return &set[K]{d: dict.Empty[K, struct{}]()}
}

// Create a set with one value.
func Singleton[K Comparable[K]](v K) Set[K] {
	return &set[K]{d: dict.Singleton(v, struct{}{})}
}

// Insert a value into a set.
func Insert[K Comparable[K]](k K, s Set[K]) Set[K] {
	return &set[K]{d: dict.Insert(k, struct{}{}, s.set_().d)}
}

// Remove a value from a set. If the value is not found, no changes are made.
func Remove[K Comparable[K]](k K, s Set[K]) Set[K] {
	return &set[K]{d: dict.Remove(k, s.set_().d)}
}

// QUERY

// Determine if a set is empty.
func IsEmpty[K Comparable[K]](s Set[K]) bool {
	return dict.IsEmpty(s.set_().d)
}

// Determine if a value is in a set.
func Member[K Comparable[K]](k K, s Set[K]) bool {
	return dict.Member(k, s.set_().d)
}

// Determine the number of elements in a set.
func Size[K Comparable[K]](s Set[K]) Int {
	return dict.Size(s.set_().d)
}

// COMBINE

// Get the union of two sets. Keep all values.
func Union[K Comparable[K]](s1 Set[K], s2 Set[K]) Set[K] {
	return &set[K]{dict.Union(s1.set_().d, s2.set_().d)}
}

// Get the intersection of two sets. Keeps values that appear in both sets.
func Intersect[K Comparable[K]](s1 Set[K], s2 Set[K]) Set[K] {
	return &set[K]{dict.Intersect(s1.set_().d, s2.set_().d)}
}

// Get the difference between the first set and the second. Keeps values that do not appear in the second set.
func Diff[K Comparable[K]](s1 Set[K], s2 Set[K]) Set[K] {
	return &set[K]{dict.Diff(s1.set_().d, s2.set_().d)}
}

// LISTS

// Convert a set into a list, sorted from lowest to highest.
func ToList[K Comparable[K]](s Set[K]) list.List[K] {
	return dict.Keys(s.set_().d)
}

// Convert a list into a set, removing any duplicates.
func FromList[K Comparable[K]](xs list.List[K]) Set[K] {
	return list.Foldl[K, Set[K]](Insert, Empty[K](), xs)
}

// TRANSFORM

// Map a function onto a set, creating a new set with no duplicates.
func Map[A Comparable[A], B Comparable[B]](f func(A) B, s Set[A]) Set[B] {
	return FromList(
		Foldl(
			func(x A, s list.List[B]) list.List[B] { return list.Cons(f(x), s) },
			list.Empty[B](),
			s,
		),
	)
}

// Fold over the values in a set, in order from lowest to highest.
func Foldl[A Comparable[A], B any](f func(A, B) B, initialState B, s Set[A]) B {
	return dict.Foldl(func(key A, _ struct{}, state B) B {
		return f(key, state)
	},
		initialState,
		s.set_().d,
	)
}

// Fold over the values in a set, in order from highest to lowest.
func Foldr[A Comparable[A], B any](f func(A, B) B, initialState B, s Set[A]) B {
	return dict.Foldr(func(key A, _ struct{}, state B) B {
		return f(key, state)
	},
		initialState,
		s.set_().d,
	)
}

// Only keep elements that pass the given test.
func Filter[A Comparable[A]](isGood func(A) bool, s Set[A]) Set[A] {
	return &set[A]{dict.Filter(func(k A, _ struct{}) bool { return isGood(k) }, s.set_().d)}
}

// Create two new sets. The first contains all the elements that passed the
// given test, and the second contains all the elements that did not.
func Partition[A Comparable[A]](isGood func(A) bool, s Set[A]) tuple.Tuple2[Set[A], Set[A]] {
	_v0 := dict.Partition(func(k A, _ struct{}) bool { return isGood(k) }, s.set_().d)
	dict1 := tuple.First(_v0)
	dict2 := tuple.Second(_v0)
	return tuple.Pair[Set[A], Set[A]](&set[A]{dict1}, &set[A]{dict2})
}
