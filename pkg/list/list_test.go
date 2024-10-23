package list

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	kernel "github.com/Confidenceman02/scion-tools/pkg/list/internal"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmp(t *testing.T) {
	asserts := assert.New(t)

	t.Run("When empty equal", func(t *testing.T) {
		l1 := Empty[Int]()
		l2 := Empty[Int]()

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(EQ{}, Compare(l1, l2))
	})

	t.Run("When cons equal Int", func(t *testing.T) {
		l1 := Singleton[Int](1)
		l2 := Singleton[Int](1)
		l3 := Range(1, 10)
		l4 := Range(1, 10)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(0, l3.Cmp(l4))
		asserts.Equal(EQ{}, Compare(l1, l2))
		asserts.Equal(EQ{}, Compare(l3, l4))
	})

	t.Run("When cons equal Float", func(t *testing.T) {
		l1 := Singleton[Float](1.0)
		l2 := Singleton[Float](1.0)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(EQ{}, Compare(l1, l2))
	})

	t.Run("When cons greater Int", func(t *testing.T) {
		l1 := Singleton[Int](2)
		l2 := Singleton[Int](1)
		l3 := Range(1, 11)
		l4 := Range(1, 10)

		asserts.Equal(+1, l1.Cmp(l2))
		asserts.Equal(+1, l3.Cmp(l4))
		asserts.Equal(GT{}, Compare(l1, l2))
		asserts.Equal(GT{}, Compare(l3, l4))
	})

	t.Run("When cons greater Float", func(t *testing.T) {
		l1 := Singleton[Float](1.1)
		l2 := Singleton[Float](1.0)
		l3 := fromArray([]Float{1.1, 2.2})
		l4 := fromArray([]Float{1.1, 2.1})

		asserts.Equal(+1, l1.Cmp(l2))
		asserts.Equal(+1, l3.Cmp(l4))
		asserts.Equal(GT{}, Compare(l1, l2))
		asserts.Equal(GT{}, Compare(l3, l4))
	})

	t.Run("When cons less Int", func(t *testing.T) {
		l1 := Singleton[Int](1)
		l2 := Singleton[Int](2)
		l3 := Range(1, 10)
		l4 := Range(1, 11)

		asserts.Equal(-1, l1.Cmp(l2))
		asserts.Equal(-1, l3.Cmp(l4))
		asserts.Equal(LT{}, Compare(l1, l2))
		asserts.Equal(LT{}, Compare(l3, l4))
	})

	t.Run("When cons less Float", func(t *testing.T) {
		l1 := Singleton[Float](1)
		l2 := Singleton[Float](2)
		l3 := fromArray([]Float{1.0, 2.0})
		l4 := fromArray([]Float{1.0, 2.1})

		asserts.Equal(-1, l1.Cmp(l2))
		asserts.Equal(-1, l3.Cmp(l4))
		asserts.Equal(LT{}, Compare(l1, l2))
		asserts.Equal(LT{}, Compare(l3, l4))
	})

	t.Run("When empty and cons Int", func(t *testing.T) {
		l1 := Empty[Int]()
		l2 := Singleton[Int](2)

		asserts.Equal(-1, l1.Cmp(l2))
		asserts.Equal(LT{}, Compare(l1, l2))
	})

	t.Run("When empty and cons Float", func(t *testing.T) {
		l1 := Empty[Float]()
		l2 := Singleton[Float](2)

		asserts.Equal(-1, l1.Cmp(l2))
		asserts.Equal(LT{}, Compare(l1, l2))
	})
}

func TestCreateFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[Int](10)

		asserts.Equal(&list[Int]{&kernel.Cons[Int, List[Int]]{Head: 10, Tail: empty[Int]{}}}, SUT)
	})

	t.Run("Repeat", func(t *testing.T) {
		SUT := Repeat(2, 10)

		asserts.Equal(
			&list[int]{
				&kernel.Cons[int, List[int]]{
					Head: 10,
					Tail: &list[int]{
						&kernel.Cons[int, List[int]]{
							Head: 10,
							Tail: empty[int]{},
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
				&kernel.Cons[Int, List[Int]]{
					Head: 2,
					Tail: &list[Int]{
						&kernel.Cons[Int, List[Int]]{
							Head: 3,
							Tail: empty[Int]{},
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
				&kernel.Cons[int, List[int]]{
					Head: 20,
					Tail: ls,
				},
			},
			SUT,
		)
	})
}

func TestTransformFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map", func(t *testing.T) {
		SUT1 := fromArray([]Float{Float(1), Float(4), Float(9)})
		SUT2 := fromArray([]bool{true, false, true})

		asserts.Equal(&list[Float]{
			&kernel.Cons[Float, List[Float]]{
				Head: 1,
				Tail: &list[Float]{
					&kernel.Cons[Float, List[Float]]{
						Head: 2,
						Tail: &list[Float]{
							&kernel.Cons[Float, List[Float]]{
								Head: 3,
								Tail: empty[Float]{},
							},
						},
					},
				},
			},
		},
			Map(Sqrt, SUT1),
		)

		asserts.Equal(&list[bool]{
			&kernel.Cons[bool, List[bool]]{
				Head: false,
				Tail: &list[bool]{
					&kernel.Cons[bool, List[bool]]{
						Head: true,
						Tail: &list[bool]{
							&kernel.Cons[bool, List[bool]]{
								Head: false,
								Tail: empty[bool]{},
							},
						},
					},
				},
			},
		},
			Map(Not, SUT2),
		)
	})

	t.Run("Foldl", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl(Add, 0, ls)

			asserts.Equal(Int(6), SUT)
		})
		t.Run("Concat", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl[Int, List[Int]](Cons, Empty[Int](), ls)

			asserts.Equal(
				&list[Int]{
					&kernel.Cons[Int, List[Int]]{
						Head: 3,
						Tail: &list[Int]{
							&kernel.Cons[Int, List[Int]]{
								Head: 2,
								Tail: &list[Int]{
									&kernel.Cons[Int, List[Int]]{
										Head: 1,
										Tail: empty[Int]{},
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
	t.Run("Foldr", func(t *testing.T) {
		ls1 := Range(1, 3)
		ls2 := Range(1, 3)
		SUT1 := Foldr(Add, 0, ls1)
		SUT2 := Foldr(Cons[Int], Empty[Int](), ls2)

		asserts.Equal(Int(6), SUT1)
		asserts.Equal(
			&list[Int]{
				&kernel.Cons[Int, List[Int]]{
					Head: 1,
					Tail: &list[Int]{
						&kernel.Cons[Int, List[Int]]{
							Head: 2,
							Tail: &list[Int]{
								&kernel.Cons[Int, List[Int]]{
									Head: 3,
									Tail: empty[Int]{},
								},
							},
						},
					},
				},
			},
			SUT2,
		)
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
				&kernel.Cons[Int, List[Int]]{
					Head: 3,
					Tail: &list[Int]{
						&kernel.Cons[Int, List[Int]]{
							Head: 2,
							Tail: &list[Int]{
								&kernel.Cons[Int, List[Int]]{
									Head: 1,
									Tail: empty[Int]{},
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

			asserts.True(Member(9, l))
			asserts.False(Member(11, l))
		})
		t.Run("When member is List", func(t *testing.T) {
			haystack := &list[List[Int]]{
				&kernel.Cons[List[Int], List[List[Int]]]{
					Head: &list[Int]{
						&kernel.Cons[Int, List[Int]]{
							Head: 12,
							Tail: empty[Int]{},
						},
					},
				}}

			needle := &list[Int]{&kernel.Cons[Int, List[Int]]{Head: 12, Tail: empty[Int]{}}}

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
			l := Cons(2, Singleton[Int](4)) // [2,4]

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
			asserts.True(IsEmpty(SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Singleton[Int](2)

			asserts.False(IsEmpty(SUT))
		})
	})

	t.Run("Head", func(t *testing.T) {
		t.Run("When empty", func(t *testing.T) {
			SUT := Empty[Int]()
			asserts.Equal(maybe.Nothing{}, Head(SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Singleton(23)
			asserts.Equal(maybe.Just[int]{Value: 23}, Head(SUT))
		})
	})

	t.Run("Tail", func(t *testing.T) {
		t.Run("When emtpy", func(t *testing.T) {
			SUT := Empty[Int]()

			asserts.Equal(maybe.Nothing{}, Tail(SUT))
		})
		t.Run("When has empty cons", func(t *testing.T) {
			SUT := Singleton(23)

			asserts.Equal(maybe.Just[List[int]]{Value: empty[int]{}}, Tail(SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Cons(22, Singleton(23))

			asserts.Equal(
				maybe.Just[List[int]]{
					Value: &list[int]{
						&kernel.Cons[int, List[int]]{
							Head: 23,
							Tail: empty[int]{},
						},
					},
				},
				Tail(SUT),
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
