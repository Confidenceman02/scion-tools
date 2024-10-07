package list

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	. "github.com/Confidenceman02/scion-tools/pkg/maybe"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[Int](10)

		asserts.Equal(&list[Int]{consList{}, &cons[Int]{10, empty[Int]{}}}, SUT)
	})

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
		t.Run("hi is lower than low", func(t *testing.T) {
			SUT := Range(2, 1)
			asserts.Equal(Empty[Int](), SUT)
		})
	})

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

func TestTransformFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Foldl", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl[Int](Add, 0, ls)

			asserts.Equal(Int(6), SUT)
		})
		t.Run("Concat", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl[Int, List[Int]](Cons, Empty[Int](), ls)

			asserts.Equal(
				&list[Int]{
					consList{},
					&cons[Int]{
						head: 3,
						tail: &list[Int]{
							consList{},
							&cons[Int]{
								head: 2,
								tail: &list[Int]{
									consList{},
									&cons[Int]{
										head: 1,
										tail: empty[Int]{},
									},
								},
							},
						},
					},
				},
				SUT,
			)
		})
	})
}

func TestUtilityFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Length", func(t *testing.T) {
		t.Run("Empty list", func(t *testing.T) {
			ls := Empty[Int]()

			asserts.Equal(Int(0), Length(ls))
		})
		t.Run("With cons", func(t *testing.T) {
			ls := Range(1, 3)

			asserts.Equal(Int(3), Length(ls))
		})
	})

	t.Run("Reverse", func(t *testing.T) {
		ls := Range(1, 3)
		SUT := Reverse[Int](ls)

		asserts.Equal(
			&list[Int]{
				consList{},
				&cons[Int]{
					head: 3,
					tail: &list[Int]{
						consList{},
						&cons[Int]{
							head: 2,
							tail: &list[Int]{
								consList{},
								&cons[Int]{
									head: 1,
									tail: empty[Int]{},
								},
							},
						},
					},
				},
			},
			SUT,
		)

	})

	t.Run("Any", func(t *testing.T) {
		t.Run("When true", func(t *testing.T) {
			ls := Range(1, 10)

			asserts.True(Any(func(a Int) bool { return a >= 10 }, ls))
		})
		t.Run("When false", func(t *testing.T) {
			ls := Range(1, 10)

			asserts.False(Any(func(a Int) bool { return a > 10 }, ls))
		})
	})
}

func TestDeconstructFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("When empty", func(t *testing.T) {
			SUT := Empty[Int]()
			asserts.True(IsEmpty[Int](SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Singleton[Int](2)

			asserts.False(IsEmpty[Int](SUT))
		})
	})

	t.Run("Head", func(t *testing.T) {
		t.Run("When empty", func(t *testing.T) {
			SUT := Empty[Int]()
			asserts.Equal(Nothing{}, Head[Int](SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Singleton(23)
			asserts.Equal(Just[int]{Value: 23}, Head[int](SUT))
		})
	})

	t.Run("Tail", func(t *testing.T) {
		t.Run("When emtpy", func(t *testing.T) {
			SUT := Empty[Int]()

			asserts.Equal(Nothing{}, Tail[Int](SUT))
		})
		t.Run("When has empty cons", func(t *testing.T) {
			SUT := Singleton(23)

			asserts.Equal(Just[List[int]]{Value: empty[int]{}}, Tail[int](SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Cons(22, Singleton(23))

			asserts.Equal(
				Just[List[int]]{
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

	})
}

func TestEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		SUT := Empty[Int]()

		asserts.Equal(empty[Int]{}, SUT)
	})
}
