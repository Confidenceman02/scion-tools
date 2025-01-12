package result

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	m "github.com/Confidenceman02/scion-tools/pkg/maybe"
	s "github.com/Confidenceman02/scion-tools/pkg/string"
	"github.com/stretchr/testify/assert"
	"testing"
)

func add3(a Int, b Int, c Int) Int {
	return Add(a, Add(b, c))
}
func add4(a Int, b Int, c Int, d Int) Int {
	return Add(a, Add(b, Add(c, d)))
}
func add5(a Int, b Int, c Int, d Int, e Int) Int {
	return Add(a, Add(b, Add(c, Add(d, e))))
}
func isEven(n Int) Result[s.String, Int] {
	if ModBy(2, n) == 0 {
		return Ok[s.String, Int]{Val: n}
	} else {
		return Err[s.String, Int]{Err: "number is odd"}
	}
}

func toIntResult(strInt s.String) Result[s.String, Int] {
	return m.MaybeWith(
		s.ToInt(strInt),
		func(j m.Just[Int]) Result[s.String, Int] {
			return Ok[s.String, Int]{Val: j.Value}
		},
		func(n m.Nothing) Result[s.String, Int] {
			return Err[s.String, Int]{Err: "could not convert '" + strInt + "'" + " to an Int"}
		},
	)
}

func TestMap(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map Ok", func(t *testing.T) {
		asserts.Equal(
			Ok[struct{}, Int]{Val: Int(3)},
			Map(func(a Int) Int { return Add(a, 1) }, Ok[struct{}, Int]{Val: Int(2)}),
		)
	})
	t.Run("Map Err", func(t *testing.T) {
		asserts.Equal(
			Err[string, Int]{Err: "error"},
			Map(func(a Int) Int { return Add(a, 1) }, Err[string, Int]{Err: "error"}),
		)
	})
}

func TestMapN(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map2 Ok", func(t *testing.T) {
		asserts.Equal(Ok[struct{}, Int]{Val: 3}, Map2(Add, Ok[struct{}, Int]{Val: 2}, Ok[struct{}, Int]{Val: 1}))
	})
	t.Run("Map2 Err", func(t *testing.T) {
		asserts.Equal(Err[string, Int]{Err: "x"}, Map2(Add, Ok[string, Int]{Val: 1}, Err[string, Int]{Err: "x"}))
	})
	t.Run("Map3 Ok", func(t *testing.T) {
		asserts.Equal(
			Ok[struct{}, Int]{Val: 6},
			Map3(add3, Ok[struct{}, Int]{Val: 1}, Ok[struct{}, Int]{Val: 2}, Ok[struct{}, Int]{Val: 3}),
		)
	})
	t.Run("Map3 Err", func(t *testing.T) {
		asserts.Equal(
			Err[struct{}, Int]{Err: struct{}{}},
			Map3(add3, Ok[struct{}, Int]{Val: 1}, Ok[struct{}, Int]{Val: 2}, Err[struct{}, Int]{Err: struct{}{}}),
		)
	})
	t.Run("Map4 Ok", func(t *testing.T) {
		asserts.Equal(
			Ok[struct{}, Int]{Val: 10},
			Map4(add4, Ok[struct{}, Int]{Val: 1}, Ok[struct{}, Int]{Val: 2}, Ok[struct{}, Int]{Val: 3}, Ok[struct{}, Int]{Val: 4}),
		)
	})
	t.Run("Map4 Err", func(t *testing.T) {
		asserts.Equal(
			Err[struct{}, Int]{Err: struct{}{}},
			Map4(add4, Ok[struct{}, Int]{Val: 1}, Ok[struct{}, Int]{Val: 2}, Ok[struct{}, Int]{Val: 3}, Err[struct{}, Int]{Err: struct{}{}}),
		)
	})
	t.Run("Map5 Ok", func(t *testing.T) {
		asserts.Equal(
			Ok[struct{}, Int]{Val: 15},
			Map5(add5, Ok[struct{}, Int]{Val: 1}, Ok[struct{}, Int]{Val: 2}, Ok[struct{}, Int]{Val: 3}, Ok[struct{}, Int]{Val: 4}, Ok[struct{}, Int]{Val: 5}),
		)
	})
	t.Run("Map5 Err", func(t *testing.T) {
		asserts.Equal(
			Err[struct{}, Int]{Err: struct{}{}},
			Map5(add5, Ok[struct{}, Int]{Val: 1}, Ok[struct{}, Int]{Val: 2}, Ok[struct{}, Int]{Val: 3}, Ok[struct{}, Int]{Val: 4}, Err[struct{}, Int]{Err: struct{}{}}),
		)
	})
}

func TestAndThen(t *testing.T) {
	asserts := assert.New(t)
	t.Run("AndThen Ok", func(t *testing.T) {
		asserts.Equal(Ok[s.String, Int]{Val: 42}, AndThen(isEven, toIntResult("42")))
	})
	t.Run("AndThen first Err", func(t *testing.T) {
		asserts.Equal(
			Err[s.String, Int]{Err: "could not convert '4.2' to an Int"},
			AndThen(isEven, toIntResult("4.2")),
		)
	})
	t.Run("AndThen second Err", func(t *testing.T) {
		asserts.Equal(
			Err[s.String, Int]{Err: "number is odd"},
			AndThen(isEven, toIntResult("41")),
		)
	})
}
