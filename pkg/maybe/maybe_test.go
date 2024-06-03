package maybe

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map a Nothing", func(t *testing.T) {
		var SUT Maybe[int]
		SUT = Nothing{}

		asserts.Equal(Nothing{}, Map(func(a int) int { return 2 }, SUT))
	})

	t.Run("Map a Just", func(t *testing.T) {
		var SUT Maybe[int]
		SUT = Just[int]{Value: 22}

		asserts.Equal(Just[string]{Value: "Hello"},
			Map(func(_ int) string { return "Hello" }, SUT),
		)
	})

	t.Run("Map2 a Nothing", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]

		m1 = Just[int]{Value: 22}
		m2 = Nothing{}

		asserts.Equal(Nothing{},
			Map2(basics.Add[int], m1, m2),
		)
	})

	t.Run("Map2 a Just", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]

		m1 = Just[int]{Value: 20}
		m2 = Just[int]{Value: 20}

		asserts.Equal(Just[int]{Value: 40},
			Map2(basics.Add[int], m1, m2),
		)
	})
}

func TestAndThen(t *testing.T) {
	asserts := assert.New(t)

	t.Run("AndThen a Nothing", func(t *testing.T) {
		var SUT Maybe[int]

		SUT = Nothing{}

		asserts.Equal(
			Nothing{},
			AndThen(func(_ int) Maybe[string] { return Just[string]{Value: "Hello"} }, SUT),
		)
	})

	t.Run("AndThen a Just", func(t *testing.T) {
		var SUT Maybe[int]

		SUT = Just[int]{Value: 22}

		asserts.Equal(
			Just[string]{Value: "Hello"},
			AndThen(func(_ int) Maybe[string] { return Just[string]{Value: "Hello"} }, SUT),
		)
	})
}

func TestWithDefault(t *testing.T) {
	asserts := assert.New(t)

	t.Run("WithDefault with Nothing", func(t *testing.T) {
		var SUT Maybe[int]

		SUT = Nothing{}

		asserts.Equal(22, WithDefault(22, SUT))
	})

	t.Run("WithDefault with Just", func(t *testing.T) {
		var SUT Maybe[int]

		SUT = Just[int]{Value: 2}

		asserts.Equal(2, WithDefault(22, SUT))
	})
}
