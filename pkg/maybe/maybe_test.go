package maybe

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
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
		var m1 Maybe[Int]
		var m2 Maybe[Int]

		m1 = Just[Int]{Value: 22}
		m2 = Nothing{}

		asserts.Equal(Nothing{},
			Map2(Add[Int], m1, m2),
		)
	})

	t.Run("Map2 a Just", func(t *testing.T) {
		var m1 Maybe[Int]
		var m2 Maybe[Int]

		m1 = Just[Int]{Value: 20}
		m2 = Just[Int]{Value: 20}

		asserts.Equal(Just[Int]{Value: 40},
			Map2(Add[Int], m1, m2),
		)
	})

	t.Run("Map3 a Nothing", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]
		var m3 Maybe[int]

		m1 = Just[int]{Value: 22}
		m2 = Just[int]{Value: 23}
		m3 = Nothing{}

		asserts.Equal(Nothing{},
			Map3(func(_ int, _ int, _ int) int { return 1 }, m1, m2, m3),
		)
	})

	t.Run("Map3 a Just", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]
		var m3 Maybe[int]

		m1 = Just[Int]{Value: 1}
		m2 = Just[Int]{Value: 2}
		m3 = Just[Int]{Value: 3}

		asserts.Equal(Just[Int]{Value: 6},
			Map3(func(a Int, b Int, c Int) Int { return Add(Add(a, b), c) }, m1, m2, m3),
		)
	})

	t.Run("Map4 a Nothing", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]
		var m3 Maybe[int]
		var m4 Maybe[int]

		m1 = Just[int]{Value: 22}
		m2 = Just[int]{Value: 23}
		m3 = Just[int]{Value: 23}
		m4 = Nothing{}

		asserts.Equal(Nothing{},
			Map4(func(_ int, _ int, _ int, _ int) int { return 1 }, m1, m2, m3, m4),
		)
	})

	t.Run("Map4 a Just", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]
		var m3 Maybe[int]
		var m4 Maybe[int]

		m1 = Just[int]{Value: 1}
		m2 = Just[int]{Value: 2}
		m3 = Just[int]{Value: 3}
		m4 = Just[int]{Value: 3}

		asserts.Equal(Just[int]{Value: 22},
			Map4(func(a int, b int, c int, d int) int { return 22 }, m1, m2, m3, m4),
		)
	})

	t.Run("Map5 a Nothing", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]
		var m3 Maybe[int]
		var m4 Maybe[int]
		var m5 Maybe[int]

		m1 = Just[int]{Value: 22}
		m2 = Just[int]{Value: 23}
		m3 = Just[int]{Value: 23}
		m4 = Just[int]{Value: 23}
		m5 = Nothing{}

		asserts.Equal(Nothing{},
			Map5(func(_ int, _ int, _ int, _ int, _ int) int { return 1 }, m1, m2, m3, m4, m5),
		)
	})

	t.Run("Map5 a Just", func(t *testing.T) {
		var m1 Maybe[int]
		var m2 Maybe[int]
		var m3 Maybe[int]
		var m4 Maybe[int]
		var m5 Maybe[int]

		m1 = Just[int]{Value: 1}
		m2 = Just[int]{Value: 2}
		m3 = Just[int]{Value: 3}
		m4 = Just[int]{Value: 3}
		m5 = Just[int]{Value: 3}

		asserts.Equal(Just[int]{Value: 22},
			Map5(func(a int, b int, c int, d int, e int) int { return 22 }, m1, m2, m3, m4, m5),
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
