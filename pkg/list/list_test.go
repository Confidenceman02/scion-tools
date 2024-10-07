package list

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleton(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[Int](10)

		asserts.Equal(&list[Int]{consList{}, &cons[Int]{10, empty[Int]{}}}, SUT)
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
			&list[Int]{
				consList{},
				&cons[Int]{
					head: 2,
					tail: &list[Int]{consList{},
						&cons[Int]{
							head: 3,
							tail: empty[Int]{},
						},
					},
				},
			},
			SUT,
		)
	})

	t.Run("Range - hi is lower than low", func(t *testing.T) {
		SUT := Range(2, 1)
		asserts.Equal(Empty[Int](), SUT)
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
		SUT := Empty[Int]()
		asserts.True(IsEmpty[Int](SUT))
	})

	t.Run("When has cons", func(t *testing.T) {
		SUT := Singleton[Int](2)

		asserts.False(IsEmpty[Int](SUT))
	})
}

func TestHead(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When empty", func(t *testing.T) {
		SUT := Empty[Int]()
		asserts.Equal(maybe.Nothing{}, Head[Int](SUT))
	})

	t.Run("When has cons", func(t *testing.T) {
		SUT := Singleton(23)
		asserts.Equal(maybe.Just[int]{Value: 23}, Head[int](SUT))
	})
}

func TestTail(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When emtpy", func(t *testing.T) {
		SUT := Empty[Int]()

		asserts.Equal(maybe.Nothing{}, Tail[Int](SUT))
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
		SUT := Empty[Int]()

		asserts.Equal(empty[Int]{}, SUT)
	})
}
