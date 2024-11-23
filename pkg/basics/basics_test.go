package basics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMath(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Add", func(t *testing.T) {
		var SUT Int
		SUT = 1

		asserts.Equal(Int(2), Add(SUT, SUT))
	})

	t.Run("Fdiv", func(t *testing.T) {
		SUT1 := Fdiv(10, 4)
		SUT2 := Fdiv(11, 4)
		SUT3 := Fdiv(12, 4)
		SUT4 := Fdiv(13, 4)
		SUT5 := Fdiv(14, 4)
		SUT6 := Fdiv(-1, 4)
		SUT7 := Fdiv(-5, 4)

		asserts.Equal(Float(2.5), SUT1)
		asserts.Equal(Float(2.75), SUT2)
		asserts.Equal(Float(3), SUT3)
		asserts.Equal(Float(3.25), SUT4)
		asserts.Equal(Float(3.5), SUT5)
		asserts.Equal(Float(-0.25), SUT6)
		asserts.Equal(Float(-1.25), SUT7)

	})
}

func TestIntToFloatFloatToInt(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ToFloat", func(t *testing.T) {
		var SUT Int
		SUT = 23

		asserts.Equal(Float(23), ToFloat(SUT))
	})

	t.Run("Round 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Round(SUT))
	})

	t.Run("Round 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(1), Round(SUT))
	})

	t.Run("Round 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(2), Round(SUT))
	})

	t.Run("Round 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(2), Round(SUT))
	})

	t.Run("Floor 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Floor 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Floor 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Floor 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(1), Floor(SUT))
	})

	t.Run("Ceiling 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Ceiling(SUT))
	})

	t.Run("Ceiling 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(2), Ceiling(SUT))
	})

	t.Run("Ceiling 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(2), Ceiling(SUT))
	})

	t.Run("Ceiling 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(2), Ceiling(SUT))
	})

	t.Run("Truncate 1.0", func(t *testing.T) {
		var SUT Float
		SUT = 1.0

		asserts.Equal(Int(1), Truncate(SUT))
	})

	t.Run("Truncate 1.2", func(t *testing.T) {
		var SUT Float
		SUT = 1.2

		asserts.Equal(Int(1), Truncate(SUT))
	})

	t.Run("Truncate 1.5", func(t *testing.T) {
		var SUT Float
		SUT = 1.5

		asserts.Equal(Int(1), Truncate(SUT))
	})

	t.Run("Truncate 1.8", func(t *testing.T) {
		var SUT Float
		SUT = 1.8

		asserts.Equal(Int(1), Truncate(SUT))
	})
}

func TestEquality(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Eq", func(t *testing.T) {
		t.Run("When true", func(t *testing.T) {
			SUT := 23

			asserts.True(Eq(23, SUT))
		})
		t.Run("When false", func(t *testing.T) {
			SUT := 23

			asserts.False(Eq(22, SUT))
		})
		t.Run("When function", func(t *testing.T) {
			SUT := func() int { return 34 }

			f := func() bool { return Eq(SUT, SUT) }

			asserts.Panics(func() { f() })
		})
		t.Run("When list", func(t *testing.T) {
			map_1 := map[int]string{

				200: "Anita",
				201: "Neha",
				203: "Suman",
				204: "Robin",
				205: "Rohit",
			}
			map_2 := map[int]string{

				200: "Anita",
				201: "Neha",
				203: "Suman",
				204: "Robin",
				205: "Rohit",
			}
			asserts.True(Eq(map_1, map_2))
		})
	})
}

func TestComparisons(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Lt", func(t *testing.T) {
		asserts.True(Lt(Int(1), Int(3)))
		asserts.True(Lt(Int(-9), Int(-2)))
		asserts.False(Lt(Int(2), Int(1)))
	})

	t.Run("Gt", func(t *testing.T) {
		asserts.True(Gt(Int(3), Int(1)))
		asserts.True(Gt(Int(-2), Int(-9)))
		asserts.False(Gt(Int(1), Int(2)))
	})

	t.Run("Le", func(t *testing.T) {
		asserts.True(Le(Int(3), Int(3)))
		asserts.True(Le(Int(2), Int(3)))
		asserts.False(Le(Int(2), Int(1)))
	})

	t.Run("Ge", func(t *testing.T) {
		asserts.True(Ge(Int(3), Int(3)))
		asserts.True(Ge(Int(3), Int(2)))
		asserts.False(Ge(Int(2), Int(3)))
	})

	t.Run("Max", func(t *testing.T) {
		asserts.Equal(Int(2), Max(Int(1), Int(2)))
		asserts.Equal(Int(3), Max(Int(1), Int(3)))
		asserts.Equal(Int(12345678), Max(Int(42), Int(12345678)))
	})

	t.Run("Min", func(t *testing.T) {
		asserts.Equal(Int(42), Min(Int(42), Int(12345678)))
	})
}

func TestBooleans(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Not", func(t *testing.T) {
		asserts.True(Not(false))
		asserts.False(Not(true))
	})
}

func TestFancierMath(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ModBy", func(t *testing.T) {
		asserts.Equal(Int(0), ModBy(2, 0))
		asserts.Equal(Int(1), ModBy(2, 1))
		asserts.Equal(Int(0), ModBy(2, 2))
		asserts.Equal(Int(1), ModBy(2, 3))
		asserts.Equal(Int(3), ModBy(4, -5))
		asserts.Equal(Int(0), ModBy(4, -4))
		asserts.Equal(Int(1), ModBy(4, -3))
		asserts.Equal(Int(2), ModBy(4, -2))
		asserts.Equal(Int(3), ModBy(4, -1))
		asserts.Equal(Int(0), ModBy(4, 0))
		asserts.Equal(Int(1), ModBy(4, 1))
		asserts.Equal(Int(2), ModBy(4, 2))
		asserts.Equal(Int(3), ModBy(4, 3))
		asserts.Equal(Int(0), ModBy(4, 4))
		asserts.Equal(Int(1), ModBy(4, 5))
	})

	t.Run("Negate", func(t *testing.T) {
		SUT1 := Negate(Int(42))
		SUT2 := Negate(Int(-42))
		SUT3 := Negate(Int(0))

		asserts.Equal(Int(-42), SUT1)
		asserts.Equal(Int(42), SUT2)
		asserts.Equal(Int(0), SUT3)
	})

	t.Run("Sqrt", func(t *testing.T) {
		asserts.Equal(Float(6), Sqrt(36))
	})
}

func TestFunctionHelpers(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ComposeL", func(t *testing.T) {
		isEven := func(i Float) bool { return ModBy(2, Int(i)) == 0 }

		SUT1 := ComposeL(Not, ComposeL(isEven, Sqrt))
		SUT2 := ComposeL(isEven, Sqrt)

		asserts.False(SUT1(4))
		asserts.False(SUT1(36))
		asserts.True(SUT1(2))
		asserts.True(SUT2(4))
		asserts.True(SUT2(36))
	})
}
