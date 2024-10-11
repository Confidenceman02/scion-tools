package list

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	. "github.com/Confidenceman02/scion-tools/pkg/maybe"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmp(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When empty equal", func(t *testing.T) {
		l1 := Empty[Int]()
		l2 := Empty[Int]()

		asserts.Equal(0, l1.Cmp(l2))
	})

	t.Run("When cons equal Int", func(t *testing.T) {
		l1 := Singleton[Int](1)
		l2 := Singleton[Int](1)
		l3 := Range(1, 10)
		l4 := Range(1, 10)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(0, l3.Cmp(l4))
	})

	t.Run("When cons equal Float", func(t *testing.T) {
		l1 := Singleton[Float](1.0)
		l2 := Singleton[Float](1.0)

		asserts.Equal(0, l1.Cmp(l2))
	})

	t.Run("When cons greater Int", func(t *testing.T) {
		l1 := Singleton[Int](2)
		l2 := Singleton[Int](1)
		l3 := Range(1, 11)
		l4 := Range(1, 10)

		asserts.Equal(+1, l1.Cmp(l2))
		asserts.Equal(+1, l3.Cmp(l4))
	})

	t.Run("When cons greater Float", func(t *testing.T) {
		l1 := Singleton[Float](1.1)
		l2 := Singleton[Float](1.0)

		asserts.Equal(+1, l1.Cmp(l2))
	})

	t.Run("When cons less Int", func(t *testing.T) {
		l1 := Singleton[Int](1)
		l2 := Singleton[Int](2)
		l3 := Range(1, 10)
		l4 := Range(1, 11)

		asserts.Equal(-1, l1.Cmp(l2))
		asserts.Equal(-1, l3.Cmp(l4))
	})

	t.Run("When cons less Float", func(t *testing.T) {
		l1 := Singleton[Float](1)
		l2 := Singleton[Float](2)

		asserts.Equal(-1, l1.Cmp(l2))
	})

	t.Run("When empty and cons Int", func(t *testing.T) {
		l1 := Empty[Int]()
		l2 := Singleton[Int](2)

		asserts.Equal(-1, l1.Cmp(l2))
	})

	t.Run("When empty and cons Float", func(t *testing.T) {
		l1 := Empty[Float]()
		l2 := Singleton[Float](2)

		asserts.Equal(-1, l1.Cmp(l2))
	})
}

func TestCreateFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("FromArray", func(t *testing.T) {
		arr := []Int{1, 2, 3, 4}

		SUT := FromArray(arr)

		asserts.Equal(
			&list[Int]{
				consList{},
				&cons[Int]{
					head: 1,
					tail: &list[Int]{consList{},
						&cons[Int]{
							head: 2,
							tail: &list[Int]{
								consList{},
								&cons[Int]{
									head: 3,
									tail: &list[Int]{
										consList{},
										&cons[Int]{
											head: 4,
											tail: empty[Int]{},
										},
									},
								},
							},
						},
					},
				},
			},
			SUT,
		)
	})

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
		SUT := Reverse(ls)

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

	t.Run("Member", func(t *testing.T) {
		t.Run("When member is int", func(t *testing.T) {
			l := Range(1, 10)

			asserts.True(Member[Int](9, l))
			asserts.False(Member[Int](11, l))
		})
		t.Run("When member is List", func(t *testing.T) {
			haystack := &list[List[Int]]{
				consList{},
				&cons[List[Int]]{
					head: &list[Int]{
						consList{},
						&cons[Int]{
							head: 12,
							tail: empty[Int]{},
						},
					},
				}}

			needle := &list[Int]{consList{}, &cons[Int]{head: 12, tail: empty[Int]{}}}

			asserts.True(Member[List[Int]](needle, haystack))
		})
	})

	t.Run("All", func(t *testing.T) {
		t.Run("When false", func(t *testing.T) {
			isEven := func(i Int) bool { return ModBy(2, i) == 0 }
			l := Range(2, 3) // [2,3]

			asserts.False(All(isEven, l))
		})
		t.Run("When true", func(t *testing.T) {
			isEven := func(i Int) bool { return ModBy(2, i) == 0 }
			l := Cons[Int](2, Singleton[Int](4)) // [2,4]

			asserts.True(All(isEven, l))

		})
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
