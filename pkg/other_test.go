package pkg

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/string"
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
		xs := list.FromSlice([]basics.Comparable[string.String]{string.String("chuck"), string.String("alice"), string.String("bob")})
		SUT := list.Sort(xs)

		asserts.Equal([]basics.Comparable[string.String]{string.String("alice"), string.String("bob"), string.String("chuck")}, list.ToSlice(SUT))
	})

	t.Run("Sort_UNSASFE", func(t *testing.T) {
		xs := list.FromSlice([]string.String{"chuck", "alice", "bob"})
		SUT := list.Sort_UNSAFE(xs)

		asserts.Equal([]string.String{"alice", "bob", "chuck"}, list.ToSlice(SUT))
	})

	t.Run("SortBy", func(t *testing.T) {
		xs := list.FromSlice([]string.String{"chuck", "alice", "bob"})
		SUT := list.SortBy(func(s string.String) basics.Comparable[string.String] { return s }, xs)

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
