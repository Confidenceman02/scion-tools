package list

import (
	"testing"

	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[basics.Int](10)

		asserts.Equal(&list[basics.Int]{consList{}, &cons[basics.Int]{10, empty[basics.Int]{}}}, SUT)
	})
}

func TestRepeat(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Repeat", func(t *testing.T) {
		SUT := Repeat(2, 10)

		asserts.Equal(
			&list[int]{
				consList{},
				&cons[int]{
					head: 10,
					tail: &list[int]{consList{},
						&cons[int]{
							head: 10,
							tail: empty[int]{},
						},
					},
				},
			},
			SUT,
		)
	})
}

func TestRange(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Range", func(t *testing.T) {
		SUT := Range(2, 3)

		asserts.Equal(
			&list[basics.Int]{
				consList{},
				&cons[basics.Int]{
					head: 2,
					tail: &list[basics.Int]{consList{},
						&cons[basics.Int]{
							head: 3,
							tail: empty[basics.Int]{},
						},
					},
				},
			},
			SUT,
		)
	})

	t.Run("Range - hi is lower than low", func(t *testing.T) {
		SUT := Range(2, 1)
		asserts.Equal(Empty[basics.Int](), SUT)
	})
}

func TestCons(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Cons", func(t *testing.T) {
		ls := Singleton(10)
		SUT := Cons(20, ls)

		asserts.Equal(
			&list[int]{
				consList{},
				&cons[int]{
					head: 20,
					tail: ls,
				},
			},
			SUT,
		)
	})
}

func TestIsEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When empty", func(t *testing.T) {
		SUT := Empty[basics.Int]()
		asserts.True(IsEmpty[basics.Int](SUT))
	})

	t.Run("When has cons", func(t *testing.T) {
		SUT := Singleton[basics.Int](2)

		asserts.False(IsEmpty[basics.Int](SUT))
	})
}

func TestHead(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When empty", func(t *testing.T) {
		SUT := Empty[basics.Int]()
		asserts.Equal(maybe.Nothing{}, Head[basics.Int](SUT))
	})

	t.Run("When has cons", func(t *testing.T) {
		SUT := Singleton(23)
		asserts.Equal(maybe.Just[int]{Value: 23}, Head[int](SUT))
	})
}

func TestTail(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When emtpy", func(t *testing.T) {
		SUT := Empty[basics.Int]()

		asserts.Equal(maybe.Nothing{}, Tail[basics.Int](SUT))
	})

	t.Run("When has empty cons", func(t *testing.T) {
		SUT := Singleton(23)

		asserts.Equal(maybe.Just[List[int]]{Value: empty[int]{}}, Tail[int](SUT))
	})

	t.Run("When has cons", func(t *testing.T) {
		SUT := Cons(22, Singleton(23))

		asserts.Equal(
			maybe.Just[List[int]]{
				Value: &list[int]{
					consList{},
					&cons[int]{
						23,
						empty[int]{},
					},
				},
			},
			Tail[int](SUT),
		)
	})
}

func TestEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		SUT := Empty[basics.Int]()

		asserts.Equal(empty[basics.Int]{}, SUT)
	})
}
