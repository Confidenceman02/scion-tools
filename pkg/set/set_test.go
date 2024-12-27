package set

import (
	"testing"

	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/dict"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
	"github.com/stretchr/testify/assert"
)

var set42 = Singleton(basics.Int(42))
var set1To50 = FromList(list.Range(1, 50))
var set1To100 = FromList(list.Range(1, 100))
var set51To150 = FromList(list.Range(51, 150))
var set51To100 = FromList(list.Range(51, 100))

func isLessThan51(i basics.Int) bool {
	return basics.Lt(i, 51)
}

func TestBuildFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		t.Run("returns an empty set", func(t *testing.T) {
			asserts.Equal(basics.Int(0), Size(Empty[basics.Int]()))
		})
	})
	t.Run("Singleton", func(t *testing.T) {
		t.Run("returns set with one element", func(t *testing.T) {
			asserts.Equal(basics.Int(1), Size(Singleton[basics.Int](1)))
		})
		t.Run("contains given element", func(t *testing.T) {
			asserts.True(Member(1, Singleton[basics.Int](1)))
		})
	})

	t.Run("Insert", func(t *testing.T) {
		t.Run("adds new element to empty set", func(t *testing.T) {
			asserts.Equal(set42, Insert(42, Empty[basics.Int]()))
		})
		t.Run("adds new element to a set of 100", func(t *testing.T) {
			asserts.Equal(FromList(list.Range(1, 101)), Insert(101, set1To100))
		})
		t.Run("leaves existing element intact if it contains a given element", func(t *testing.T) {
			asserts.Equal(set42, Insert(42, set42))
		})
		t.Run("leaves set of 100 intact if it contains given element", func(t *testing.T) {
			asserts.Equal(set1To100, Insert(42, set1To100))
		})
	})
	t.Run("Remove", func(t *testing.T) {
		t.Run("removes element from singleton set", func(t *testing.T) {
			asserts.Equal(Empty[basics.Int](), Remove(42, set42))
		})
		t.Run("removes element from set of 100", func(t *testing.T) {
			asserts.Equal(list.Range(1, 99), dict.Keys(Remove(100, set1To100).set_().d))
		})
		t.Run("leaves singleton set intact if it doesn't contain given element", func(t *testing.T) {
			asserts.Equal(set42, Remove(-1, set42))
		})
		t.Run("leaves set of 100 intact if it doesn't contain given element", func(t *testing.T) {
			asserts.Equal(set1To100, Remove(-1, set1To100))
		})
	})
}

func TestQueryFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("return true for empty set", func(t *testing.T) {
			asserts.True(IsEmpty(Empty[basics.Int]()))
		})
		t.Run("returns false for singleton set", func(t *testing.T) {})
		asserts.False(IsEmpty(set42))
		t.Run("returns false for set of 100", func(t *testing.T) {
			asserts.False(IsEmpty(set1To100))
		})
	})
	t.Run("Member", func(t *testing.T) {
		t.Run("returns true when given element is inside set", func(t *testing.T) {
			asserts.True(Member(42, set42))
		})
		t.Run("return true when given element inside set of 100", func(t *testing.T) {
			asserts.True(Member(42, set1To100))
		})
		t.Run("returns false for element not in singleton", func(t *testing.T) {
			asserts.False(Member(-1, set42))
		})
		t.Run("returns false for element not in set of 100", func(t *testing.T) {
			asserts.False(Member(-1, set1To100))
		})
	})
	t.Run("Size", func(t *testing.T) {
		t.Run("returns 0 for empty set", func(t *testing.T) {
			asserts.Equal(basics.Int(0), Size(Empty[basics.Int]()))
		})
		t.Run("returns 1 for singleton set", func(t *testing.T) {
			asserts.Equal(basics.Int(1), Size(set42))
		})
		t.Run("returns 100 for set of 100", func(t *testing.T) {
			asserts.Equal(basics.Int(100), Size(set1To100))
		})
	})
}

func TestCombineFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Union", func(t *testing.T) {
		t.Run("with empty set doesn't change anything", func(t *testing.T) {
			asserts.Equal(set42, Union(set42, Empty[basics.Int]()))
		})
		t.Run("with itself doesn't change anything", func(t *testing.T) {
			asserts.Equal(dict.Keys(set1To100.set_().d), dict.Keys(Union(set1To100, set1To100).set_().d))
		})
		t.Run("with subset doesn't change anything", func(t *testing.T) {
			asserts.Equal(set1To100, Union(set1To100, set42))
		})
		t.Run("with superset returns superset", func(t *testing.T) {
			asserts.Equal(set1To100, Union(set42, set1To100))
		})
		t.Run("contains elements of both singletons", func(t *testing.T) {
			asserts.Equal(list.FromSlice([]basics.Int{1, 42}), dict.Keys(Union(set42, Singleton[basics.Int](1)).set_().d))
		})
		t.Run("consists of elements from either set", func(t *testing.T) {
			asserts.Equal(dict.Keys(FromList(list.Range(1, 150)).set_().d), dict.Keys(Union(set1To100, set51To150).set_().d))
		})
	})
	t.Run("Intersect", func(t *testing.T) {
		t.Run("with empty set returns empty set", func(t *testing.T) {
			asserts.Equal(Empty[basics.Int](), Intersect(set42, Empty[basics.Int]()))
		})
		t.Run("with itself doesn't change anything", func(t *testing.T) {
			asserts.Equal(set1To100, Intersect(set1To100, set1To100))
		})
		t.Run("with subset returns subset", func(t *testing.T) {
			asserts.Equal(set42, Intersect(set1To100, set42))
		})
		t.Run("with superset doesn't change anything", func(t *testing.T) {
			asserts.Equal(set42, Intersect(set42, set1To100))
		})
		t.Run("returns empty set given disjunctive sets", func(t *testing.T) {
			asserts.Equal(Empty[basics.Int](), Intersect(set42, Singleton[basics.Int](1)))
		})
		t.Run("consists of common elements only", func(t *testing.T) {
			asserts.Equal(set51To100, Intersect(set1To100, set51To150))
		})
	})
	t.Run("Diff", func(t *testing.T) {
		t.Run("with empty set doesn't change anything", func(t *testing.T) {
			asserts.Equal(set42, Diff(set42, Empty[basics.Int]()))
		})
		t.Run("with itself returns empty set", func(t *testing.T) {
			asserts.Equal(Empty[basics.Int](), Diff(set1To100, set1To100))
		})
		t.Run("with subset returns set without subset", func(t *testing.T) {
			asserts.Equal(Remove(42, set1To100), Diff(set1To100, set42))
		})
		t.Run("with superset returns empty set", func(t *testing.T) {
			asserts.Equal(Empty[basics.Int](), Diff(set42, set1To100))
		})
		t.Run("doesn't change anything given disjunctive sets", func(t *testing.T) {
			asserts.Equal(set42, Diff(set42, Singleton[basics.Int](1)))
		})
		t.Run("only keeps values that don't appear in the second set", func(t *testing.T) {
			asserts.Equal(dict.Keys(set1To50.set_().d), dict.Keys(Diff(set1To100, set51To150).set_().d))
		})
	})
}

func TestListFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("FromList", func(t *testing.T) {
		t.Run("return empty set for empty list", func(t *testing.T) {
			asserts.Equal(Empty[basics.Int](), FromList(list.Empty[basics.Int]()))
		})
		t.Run("returns singleton set for singleton list", func(t *testing.T) {
			asserts.Equal(set42, FromList(list.Singleton[basics.Int](42)))
		})
		t.Run("returns set with unique list elements", func(t *testing.T) {
			asserts.Equal(set1To100, FromList(list.Cons(1, list.Range(1, 100))))
		})
	})
	t.Run("ToList", func(t *testing.T) {
		t.Run("returns empty list for empty set", func(t *testing.T) {
			asserts.Equal(list.Empty[basics.Int](), ToList(Empty[basics.Int]()))
		})
		t.Run("returns singleton list for singleton set", func(t *testing.T) {
			asserts.Equal(list.Singleton[basics.Int](42), ToList(set42))
		})
		t.Run("returns sorted list of set elements", func(t *testing.T) {
			asserts.Equal(list.Range(1, 100), ToList(set1To100))
		})
	})
}

func TestTransformFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map", func(t *testing.T) {
		t.Run("applies given function to singleton element", func(t *testing.T) {
			asserts.Equal(Singleton[basics.Int](43), Map(func(i basics.Int) basics.Int { return basics.Add(i, 1) }, set42))
		})
		t.Run("applies given function to each element", func(t *testing.T) {
			asserts.Equal(FromList(list.Range(-100, -1)), Map(basics.Negate, set1To100))
		})
	})

	t.Run("Foldl", func(t *testing.T) {
		t.Run("with insert and empty set acts as identity function", func(t *testing.T) {
			asserts.Equal(set1To100, Foldl(Insert[basics.Int], Empty[basics.Int](), set1To100))
		})
		t.Run("with counter ans zero acts as size function", func(t *testing.T) {
			asserts.Equal(basics.Int(100), Foldl(func(_, count basics.Int) basics.Int { return basics.Add(count, 1) }, 0, set1To100))
		})
		t.Run("folds set elements from lowest to highest", func(t *testing.T) {
			asserts.Equal(
				list.FromSlice([]basics.Int{3, 2, 1}),
				Foldl(func(n basics.Int, ns list.List[basics.Int]) list.List[basics.Int] { return list.Cons(n, ns) },
					list.Empty[basics.Int](),
					FromList(list.FromSlice([]basics.Int{2, 1, 3})),
				),
			)
		})
	})
	t.Run("Foldr", func(t *testing.T) {
		t.Run("with insert and empty set acts as identity function", func(t *testing.T) {
			asserts.Equal(ToList(set1To100), ToList(Foldr(Insert[basics.Int], Empty[basics.Int](), set1To100)))
		})
		t.Run("with counter and zero acts as size function", func(t *testing.T) {
			asserts.Equal(basics.Int(100), Foldr(func(_, count basics.Int) basics.Int { return basics.Add(count, 1) }, 0, set1To100))
		})
		t.Run("folds set elements from highest to lowest", func(t *testing.T) {
			asserts.Equal(
				list.FromSlice([]basics.Int{1, 2, 3}),
				Foldr(func(n basics.Int, ns list.List[basics.Int]) list.List[basics.Int] { return list.Cons(n, ns) }, list.Empty[basics.Int](), FromList(list.FromSlice([]basics.Int{2, 1, 3}))))
		})
	})
	t.Run("Filter", func(t *testing.T) {
		t.Run("with always true doesn't change anything", func(t *testing.T) {
			asserts.Equal(
				set1To100,
				Filter(func(i basics.Int) bool { return basics.Always(true, i) }, set1To100),
			)
		})
		t.Run("simple filter", func(t *testing.T) {
			asserts.Equal(
				set1To50,
				Filter(isLessThan51, set1To100),
			)
		})
	})
	t.Run("Partition", func(t *testing.T) {
		t.Run("of empty set return two empty sets", func(t *testing.T) {
			asserts.Equal(tuple.Pair(Empty[basics.Int](), Empty[basics.Int]()), Partition(isLessThan51, Empty[basics.Int]()))
		})
		t.Run("simple partition", func(t *testing.T) {
			asserts.Equal(ToList(set1To50), ToList(tuple.First(Partition(isLessThan51, set1To100))))
			asserts.Equal(ToList(set51To100), ToList(tuple.Second(Partition(isLessThan51, set1To100))))
		})
	})
}
