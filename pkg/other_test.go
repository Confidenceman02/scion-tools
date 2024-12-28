package pkg

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/string"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	asserts := assert.New(t)

	t.Run("FilterMap", func(t *testing.T) {
		xs := list.FromSlice([]string.String{"3", "hi", "12", "4th", "May"})
		SUT := list.FilterMap(string.ToInt, xs)

		asserts.Equal(
			[]basics.Int{3, 12},
			list.ToSlice(SUT),
		)
	})
	t.Run("Sort", func(t *testing.T) {
		xs := list.FromSlice([]string.String{string.String("chuck"), string.String("alice"), string.String("bob")})
		SUT := list.Sort(xs)

		asserts.Equal([]string.String{string.String("alice"), string.String("bob"), string.String("chuck")}, list.ToSlice(SUT))
	})

	t.Run("SortBy", func(t *testing.T) {
		xs := list.FromSlice([]string.String{"chuck", "alice", "bob"})
		SUT := list.SortBy(func(s string.String) string.String { return s }, xs)

		asserts.Equal([]string.String{"alice", "bob", "chuck"}, list.ToSlice(SUT))
	})

	t.Run("SortWith", func(t *testing.T) {
		xs := list.Range(1, 5)
		flippedComparison := func(a, b basics.Int) basics.Order {
			switch basics.Compare(a, b) {
			case basics.LT{}:
				return basics.GT{}
			case basics.EQ{}:
				return basics.EQ{}
			default:
				return basics.LT{}
			}
		}
		SUT := list.SortWith(flippedComparison, xs)

		asserts.Equal([]basics.Int{5, 4, 3, 2, 1}, list.ToSlice(SUT))
	})
}

func TestTuple(t *testing.T) {
	asserts := assert.New(t)

	t.Run("MapFirst", func(t *testing.T) {
		SUT1 := tuple.MapFirst(string.Reverse, tuple.Pair(string.String("stressed"), 16))
		SUT2 := tuple.MapFirst(string.Length, tuple.Pair(string.String("stressed"), 16))

		asserts.Equal(
			string.String("desserts"),
			tuple.First(SUT1),
		)
		asserts.Equal(
			basics.Int(8),
			tuple.First(SUT2),
		)
	})

	t.Run("MapSecond", func(t *testing.T) {
		SUT1 := tuple.MapSecond(basics.Sqrt, tuple.Pair("stressed", basics.Float(16)))
		SUT2 := tuple.MapSecond(basics.Negate, tuple.Pair("stressed", basics.Float(16)))

		asserts.Equal(
			basics.Float(4),
			tuple.Second(SUT1),
		)
		asserts.Equal(
			basics.Float(-16),
			tuple.Second(SUT2),
		)
	})

	t.Run("MapBoth", func(t *testing.T) {
		SUT1 := tuple.MapBoth(string.Reverse, basics.Sqrt, tuple.Pair(string.String("stressed"), basics.Float(16)))
		SUT2 := tuple.MapBoth(string.Length, basics.Negate, tuple.Pair(string.String("stressed"), basics.Float(16)))

		asserts.Equal(
			string.String("desserts"),
			tuple.First(SUT1),
		)
		asserts.Equal(
			basics.Float(4),
			tuple.Second(SUT1),
		)
		asserts.Equal(basics.Int(8), tuple.First(SUT2))
		asserts.Equal(basics.Float(-16), tuple.Second(SUT2))
	})
}

func TestString(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Append", func(t *testing.T) {
		asserts.Equal(string.String("helloworld"), basics.Append(string.String("hello"), string.String("world")))
	})
}
