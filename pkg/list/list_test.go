package list

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/internal"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	. "github.com/Confidenceman02/scion-tools/pkg/tuple"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmp(t *testing.T) {
	asserts := assert.New(t)

	t.Run(" Compare", func(t *testing.T) {
		x1 := FromSlice([]basics.Int{basics.Int(1), basics.Int(2), basics.Int(3)})
		x2 := FromSlice([]basics.Int{basics.Int(1), basics.Int(2), basics.Int(4)})

		asserts.Equal(basics.LT{}, basics.Compare(x1, x2))
	})

	t.Run("Compare: When empty", func(t *testing.T) {
		l1 := Empty[basics.Int]()
		l2 := Empty[basics.Int]()

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(basics.EQ{}, basics.Compare(l1, l2))
	})

	t.Run("Compare: When cons basics.Int", func(t *testing.T) {
		l1 := Singleton[basics.Int](1)
		l2 := Singleton[basics.Int](1)
		l3 := Range(1, 10)
		l4 := Range(1, 10)
		l5 := Singleton[basics.Int](2)
		l6 := Singleton[basics.Int](1)
		l7 := Range(1, 11)
		l8 := Range(1, 10)
		l9 := Singleton[basics.Int](1)
		l10 := Singleton[basics.Int](2)
		l11 := Range(1, 10)
		l12 := Range(1, 11)
		l13 := Empty[basics.Int]()
		l14 := Singleton[basics.Int](2)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(0, l3.Cmp(l4))
		asserts.Equal(basics.EQ{}, basics.Compare(l1, l2))
		asserts.Equal(basics.EQ{}, basics.Compare(l3, l4))
		asserts.Equal(+1, l5.Cmp(l6))
		asserts.Equal(+1, l7.Cmp(l8))
		asserts.Equal(basics.GT{}, basics.Compare(l5, l6))
		asserts.Equal(basics.GT{}, basics.Compare(l7, l8))
		asserts.Equal(-1, l9.Cmp(l10))
		asserts.Equal(-1, l11.Cmp(l12))
		asserts.Equal(basics.LT{}, basics.Compare(l9, l10))
		asserts.Equal(basics.LT{}, basics.Compare(l11, l12))
		asserts.Equal(-1, l13.Cmp(l14))
		asserts.Equal(basics.LT{}, basics.Compare(l13, l14))
	})

	t.Run("Compare: When cons basics.Float", func(t *testing.T) {
		l1 := Singleton[basics.Float](1.0)
		l2 := Singleton[basics.Float](1.0)
		l3 := Singleton[basics.Float](1.1)
		l4 := Singleton[basics.Float](1.0)
		l5 := FromSlice([]basics.Float{1.1, 2.2})
		l6 := FromSlice([]basics.Float{1.1, 2.1})
		l7 := Singleton[basics.Float](1)
		l8 := Singleton[basics.Float](2)
		l9 := FromSlice([]basics.Float{1.0, 2.0})
		l10 := FromSlice([]basics.Float{1.0, 2.1})
		l11 := Empty[basics.Float]()
		l12 := Singleton[basics.Float](2)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(basics.EQ{}, basics.Compare(l1, l2))
		asserts.Equal(+1, l3.Cmp(l4))
		asserts.Equal(+1, l5.Cmp(l6))
		asserts.Equal(basics.GT{}, basics.Compare(l3, l4))
		asserts.Equal(basics.GT{}, basics.Compare(l5, l6))
		asserts.Equal(-1, l7.Cmp(l8))
		asserts.Equal(-1, l9.Cmp(l10))
		asserts.Equal(basics.LT{}, basics.Compare(l7, l8))
		asserts.Equal(basics.LT{}, basics.Compare(l9, l10))
		asserts.Equal(-1, l11.Cmp(l12))
		asserts.Equal(basics.LT{}, basics.Compare(l11, l12))
	})
}

func TestCreateFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[basics.Int](10)

		asserts.Equal(&list[basics.Int]{&internal.Cons_[basics.Int, List[basics.Int]]{A: 10, B: empty[basics.Int]{}}}, SUT)
	})

	t.Run("Repeat", func(t *testing.T) {
		SUT := Repeat(2, 10)

		asserts.Equal(
			&list[int]{
				&internal.Cons_[int, List[int]]{
					A: 10,
					B: &list[int]{
						&internal.Cons_[int, List[int]]{
							A: 10,
							B: empty[int]{},
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
			&list[basics.Int]{
				&internal.Cons_[basics.Int, List[basics.Int]]{
					A: 2,
					B: &list[basics.Int]{
						&internal.Cons_[basics.Int, List[basics.Int]]{
							A: 3,
							B: empty[basics.Int]{},
						},
					},
				},
			},
			SUT,
		)
		t.Run("hi is lower than low", func(t *testing.T) {
			SUT := Range(2, 1)
			asserts.Equal(Empty[basics.Int](), SUT)
		})
	})

	t.Run("Cons", func(t *testing.T) {
		ls := Singleton(10)
		SUT := Cons(20, ls)

		asserts.Equal(
			&list[int]{
				&internal.Cons_[int, List[int]]{
					A: 20,
					B: ls,
				},
			},
			SUT,
		)
	})
}

func TestTransformFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map", func(t *testing.T) {
		SUT1 := FromSlice([]basics.Float{basics.Float(1), basics.Float(4), basics.Float(9)})
		SUT2 := FromSlice([]bool{true, false, true})

		asserts.Equal(&list[basics.Float]{
			&internal.Cons_[basics.Float, List[basics.Float]]{
				A: 1,
				B: &list[basics.Float]{
					&internal.Cons_[basics.Float, List[basics.Float]]{
						A: 2,
						B: &list[basics.Float]{
							&internal.Cons_[basics.Float, List[basics.Float]]{
								A: 3,
								B: empty[basics.Float]{},
							},
						},
					},
				},
			},
		},
			Map(basics.Sqrt, SUT1),
		)

		asserts.Equal(&list[bool]{
			&internal.Cons_[bool, List[bool]]{
				A: false,
				B: &list[bool]{
					&internal.Cons_[bool, List[bool]]{
						A: true,
						B: &list[bool]{
							&internal.Cons_[bool, List[bool]]{
								A: false,
								B: empty[bool]{},
							},
						},
					},
				},
			},
		},
			Map(basics.Not, SUT2),
		)
	})

	t.Run("IndexedMap", func(t *testing.T) {
		x1 := Range(0, 2)

		SUT := IndexedMap(basics.Add, x1)

		asserts.Equal(
			&list[basics.Int]{
				&internal.Cons_[basics.Int, List[basics.Int]]{
					A: basics.Int(0),
					B: &list[basics.Int]{
						&internal.Cons_[basics.Int, List[basics.Int]]{
							A: basics.Int(2),
							B: &list[basics.Int]{
								&internal.Cons_[basics.Int, List[basics.Int]]{
									A: basics.Int(4),
									B: empty[basics.Int]{},
								},
							},
						},
					},
				},
			},
			SUT,
		)
	})

	t.Run("Foldl", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl(basics.Add, 0, ls)

			asserts.Equal(basics.Int(6), SUT)
		})
		t.Run("Concat", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl[basics.Int, List[basics.Int]](Cons, Empty[basics.Int](), ls)

			asserts.Equal(
				&list[basics.Int]{
					&internal.Cons_[basics.Int, List[basics.Int]]{
						A: 3,
						B: &list[basics.Int]{
							&internal.Cons_[basics.Int, List[basics.Int]]{
								A: 2,
								B: &list[basics.Int]{
									&internal.Cons_[basics.Int, List[basics.Int]]{
										A: 1,
										B: empty[basics.Int]{},
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
		SUT1 := Foldr(basics.Add, 0, ls1)
		SUT2 := Foldr(Cons[basics.Int], Empty[basics.Int](), ls2)

		asserts.Equal(basics.Int(6), SUT1)
		asserts.Equal(
			&list[basics.Int]{
				&internal.Cons_[basics.Int, List[basics.Int]]{
					A: 1,
					B: &list[basics.Int]{
						&internal.Cons_[basics.Int, List[basics.Int]]{
							A: 2,
							B: &list[basics.Int]{
								&internal.Cons_[basics.Int, List[basics.Int]]{
									A: 3,
									B: empty[basics.Int]{},
								},
							},
						},
					},
				},
			},
			SUT2,
		)
	})
	t.Run("Filter", func(t *testing.T) {
		xs := Range(1, 6)
		isEven := func(i basics.Int) bool { return basics.ModBy(2, i) == 0 }

		SUT := Filter(isEven, xs)

		asserts.Equal(
			&list[basics.Int]{
				&internal.Cons_[basics.Int, List[basics.Int]]{
					A: basics.Int(2),
					B: &list[basics.Int]{
						&internal.Cons_[basics.Int, List[basics.Int]]{
							A: basics.Int(4),
							B: &list[basics.Int]{
								&internal.Cons_[basics.Int, List[basics.Int]]{
									A: basics.Int(6),
									B: empty[basics.Int]{},
								},
							},
						},
					},
				},
			},
			SUT,
		)
	})
}

func TestUtilityFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Length", func(t *testing.T) {
		t.Run("Empty list", func(t *testing.T) {
			ls := Empty[basics.Int]()

			asserts.Equal(basics.Int(0), Length(ls))
		})
		t.Run("With cons", func(t *testing.T) {
			ls := Range(1, 3)

			asserts.Equal(basics.Int(3), Length(ls))
		})
	})

	t.Run("Reverse", func(t *testing.T) {
		ls := Range(1, 3)
		SUT := Reverse(ls)

		asserts.Equal(
			&list[basics.Int]{
				&internal.Cons_[basics.Int, List[basics.Int]]{
					A: 3,
					B: &list[basics.Int]{
						&internal.Cons_[basics.Int, List[basics.Int]]{
							A: 2,
							B: &list[basics.Int]{
								&internal.Cons_[basics.Int, List[basics.Int]]{
									A: 1,
									B: empty[basics.Int]{},
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
			haystack := &list[List[basics.Int]]{
				&internal.Cons_[List[basics.Int], List[List[basics.Int]]]{
					A: &list[basics.Int]{
						&internal.Cons_[basics.Int, List[basics.Int]]{
							A: 12,
							B: empty[basics.Int]{},
						},
					},
				}}

			needle := &list[basics.Int]{&internal.Cons_[basics.Int, List[basics.Int]]{A: 12, B: empty[basics.Int]{}}}

			asserts.True(Member[List[basics.Int]](needle, haystack))
		})
	})

	t.Run("All", func(t *testing.T) {
		t.Run("When false", func(t *testing.T) {
			isEven := func(i basics.Int) bool { return basics.ModBy(2, i) == 0 }
			l := Range(2, 3) // [2,3]

			asserts.False(All(isEven, l))
		})
		t.Run("When true", func(t *testing.T) {
			isEven := func(i basics.Int) bool { return basics.ModBy(2, i) == 0 }
			l := Cons(2, Singleton[basics.Int](4)) // [2,4]

			asserts.True(All(isEven, l))

		})
	})

	t.Run("Any", func(t *testing.T) {
		t.Run("When true", func(t *testing.T) {
			ls := Range(1, 10)

			asserts.True(Any(func(a basics.Int) bool { return a >= 10 }, ls))
		})
		t.Run("When false", func(t *testing.T) {
			ls := Range(1, 10)

			asserts.False(Any(func(a basics.Int) bool { return a > 10 }, ls))
		})
	})
	t.Run("Maximum", func(t *testing.T) {
		xs := FromSlice([]basics.Int{basics.Int(1), basics.Int(4), basics.Int(2)})
		SUT := Maximum(xs)

		asserts.Equal(maybe.Just[basics.Int]{Value: basics.Int(4)}, SUT)
	})
	t.Run("Minimum", func(t *testing.T) {
		xs1 := FromSlice([]basics.Int{basics.Int(3), basics.Int(2), basics.Int(1)})
		SUT1 := Minimum(xs1)
		SUT2 := Minimum(Empty[basics.Int]())

		asserts.Equal(maybe.Just[basics.Int]{Value: basics.Int(1)}, SUT1)
		asserts.Equal(maybe.Nothing{}, SUT2)
	})
	t.Run("Sum", func(t *testing.T) {
		xs1 := FromSlice([]basics.Int{basics.Int(1), basics.Int(2), basics.Int(3)})
		xs2 := FromSlice([]basics.Int{basics.Int(1), basics.Int(1), basics.Int(1)})
		SUT1 := Sum(xs1)
		SUT2 := Sum(xs2)
		SUT3 := Sum(Empty[basics.Int]())

		asserts.Equal(basics.Int(6), SUT1)
		asserts.Equal(basics.Int(3), SUT2)
		asserts.Equal(basics.Int(0), SUT3)
	})
	t.Run("Product", func(t *testing.T) {
		xs1 := FromSlice([]basics.Int{basics.Int(2), basics.Int(2), basics.Int(2)})
		xs2 := FromSlice([]basics.Int{basics.Int(3), basics.Int(3), basics.Int(3)})
		SUT1 := Product(xs1)
		SUT2 := Product(xs2)
		SUT3 := Product(Empty[basics.Int]())

		asserts.Equal(basics.Int(8), SUT1)
		asserts.Equal(basics.Int(27), SUT2)
		asserts.Equal(basics.Int(1), SUT3)

	})
}

func TestCombineFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Append with empty", func(t *testing.T) {
		xs := Empty[basics.Int]()
		ys := Range(1, 3)
		SUT := Append(xs, ys)

		asserts.Equal(ys, SUT)
	})

	t.Run("Append", func(t *testing.T) {
		xs1 := Singleton[basics.Int](1)
		ys1 := Range(2, 3)
		xs2 := FromSlice([]basics.Int{basics.Int(1), basics.Int(1), basics.Int(2)})
		ys2 := FromSlice([]basics.Int{basics.Int(3), basics.Int(5), basics.Int(8)})
		SUT1 := Tail(Append(xs1, ys1))
		SUT2 := Append(xs2, ys2)

		asserts.Equal(maybe.Just[List[basics.Int]]{Value: ys1}, SUT1)
		asserts.Equal(&list[basics.Int]{&internal.Cons_[basics.Int, List[basics.Int]]{A: 1, B: empty[basics.Int]{}}}, xs1)
		asserts.Equal([]basics.Int{1, 1, 2, 3, 5, 8}, ToSlice(SUT2))
		// Structure sharing
		asserts.Equal(&list[basics.Int]{
			&internal.Cons_[basics.Int, List[basics.Int]]{
				A: 1,
				B: &list[basics.Int]{
					&internal.Cons_[basics.Int, List[basics.Int]]{
						A: 1,
						B: &list[basics.Int]{
							&internal.Cons_[basics.Int, List[basics.Int]]{A: 2, B: ys2},
						},
					},
				},
			},
		},
			SUT2,
		)
	})

	t.Run("Concat", func(t *testing.T) {
		lists := FromSlice([]List[basics.Int]{Range(1, 2), Range(3, 3), Range(4, 5)})
		SUT := Concat(lists)

		asserts.Equal([]basics.Int{1, 2, 3, 4, 5}, ToSlice(SUT))
	})

	t.Run("ConcatMap", func(t *testing.T) {
		fun := func(a basics.Int) List[basics.Int] { return Singleton(a + 1) }
		xs := Range(1, 3)
		SUT := ConcatMap(fun, xs)

		asserts.Equal([]basics.Int{2, 3, 4}, ToSlice(SUT))
		asserts.Equal([]basics.Int{1, 2, 3}, ToSlice(xs))
	})

	t.Run("basics.Intersperse", func(t *testing.T) {
		xs := FromSlice([]string{"turtles", "turtles", "turtles"})
		SUT := Intersperse("on", xs)

		asserts.Equal([]string{"turtles", "on", "turtles", "on", "turtles"}, ToSlice(SUT))
	})

	t.Run("Map2", func(t *testing.T) {
		t.Run("Lists with one member", func(t *testing.T) {
			a := Cons(1, Empty[basics.Int]())
			b := Cons(1, Empty[basics.Int]())
			SUT := Map2(basics.Add, a, b)

			asserts.Equal(
				&list[basics.Int]{
					&internal.Cons_[basics.Int, List[basics.Int]]{
						A: 2, B: Empty[basics.Int]()}},
				SUT,
			)
		})
		t.Run("When one list empty", func(t *testing.T) {
			a := Range(1, 3)
			b := Empty[basics.Int]()
			SUT := Map2(basics.Add, a, b)

			asserts.Equal(empty[basics.Int]{}, SUT)
		})
	})

	t.Run("Map3", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 3)
		xc := Range(1, 3)
		add3 := func(a basics.Int, b basics.Int, c basics.Int) basics.Int { return basics.Add(a, basics.Add(b, c)) }

		SUT := Map3(add3, xa, xb, xc)

		asserts.Equal(
			[]basics.Int{3, 6, 9},
			ToSlice(SUT),
		)
	})
	t.Run("Map3 with asymmetrical lists", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 2)
		xc := Range(1, 3)
		add3 := func(a basics.Int, b basics.Int, c basics.Int) basics.Int { return basics.Add(a, basics.Add(b, c)) }

		SUT := Map3(add3, xa, xb, xc)

		asserts.Equal(
			[]basics.Int{3, 6},
			ToSlice(SUT),
		)
	})

	t.Run("Map4", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 3)
		xc := Range(1, 3)
		xd := Range(1, 3)
		add4 := func(a basics.Int, b basics.Int, c basics.Int, d basics.Int) basics.Int {
			return basics.Add(a, basics.Add(b, (basics.Add(c, d))))
		}

		SUT := Map4(add4, xa, xb, xc, xd)

		asserts.Equal(
			[]basics.Int{4, 8, 12},
			ToSlice(SUT),
		)
	})

	t.Run("Map5", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 3)
		xc := Range(1, 3)
		xd := Range(1, 3)
		xe := Range(1, 3)
		add5 := func(a basics.Int, b basics.Int, c basics.Int, d basics.Int, e basics.Int) basics.Int {
			return basics.Add(a, basics.Add(b, (basics.Add(c, basics.Add(d, e)))))
		}

		SUT := Map5(add5, xa, xb, xc, xd, xe)

		asserts.Equal(
			[]basics.Int{5, 10, 15},
			ToSlice(SUT),
		)
	})
}

func testSortFUnctions(t *testing.T) {
	// asserts := assert.New(t)

}

func TestDeconstructFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("When empty", func(t *testing.T) {
			SUT := Empty[basics.Int]()
			asserts.True(IsEmpty(SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Singleton[basics.Int](2)

			asserts.False(IsEmpty(SUT))
		})
	})

	t.Run("Head", func(t *testing.T) {
		t.Run("When empty", func(t *testing.T) {
			SUT := Empty[basics.Int]()
			asserts.Equal(maybe.Nothing{}, Head(SUT))
		})
		t.Run("When has cons", func(t *testing.T) {
			SUT := Singleton(23)
			asserts.Equal(maybe.Just[int]{Value: 23}, Head(SUT))
		})
	})

	t.Run("Tail", func(t *testing.T) {
		t.Run("When emtpy", func(t *testing.T) {
			SUT := Empty[basics.Int]()

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
						&internal.Cons_[int, List[int]]{
							A: 23,
							B: empty[int]{},
						},
					},
				},
				Tail(SUT),
			)
		})
	})

	t.Run("Take", func(t *testing.T) {
		xs := Range(1, 5)
		SUT1 := Take(1, xs)
		SUT2 := Take(2, xs)
		SUT3 := Take(3, xs)
		SUT4 := Take(4, xs)
		SUT5 := Take(5, xs)
		SUT6 := Take(1, Range(1, 1))
		SUT7 := Take(2, Range(1, 2))
		SUT8 := Take(3, Range(1, 3))
		SUT9 := Take(4, Range(1, 4))
		SUT10 := Take(100, xs)

		asserts.Equal([]basics.Int{1}, ToSlice(SUT1))
		asserts.Equal([]basics.Int{1, 2}, ToSlice(SUT2))
		asserts.Equal([]basics.Int{1, 2, 3}, ToSlice(SUT3))
		asserts.Equal([]basics.Int{1, 2, 3, 4}, ToSlice(SUT4))
		asserts.Equal([]basics.Int{1, 2, 3, 4, 5}, ToSlice(SUT5))
		asserts.Equal([]basics.Int{1}, ToSlice(SUT6))
		asserts.Equal([]basics.Int{1, 2}, ToSlice(SUT7))
		asserts.Equal([]basics.Int{1, 2, 3}, ToSlice(SUT8))
		asserts.Equal([]basics.Int{1, 2, 3, 4}, ToSlice(SUT9))
		asserts.Equal([]basics.Int{1, 2, 3, 4, 5}, ToSlice(SUT10))
	})

	t.Run("Drop", func(t *testing.T) {
		xs := Range(1, 4)
		SUT := Drop(2, xs)

		asserts.Equal([]basics.Int{3, 4}, ToSlice(SUT))
	})

	t.Run("Partition", func(t *testing.T) {
		xs := Range(1, 5)
		SUT := Partition(func(x basics.Int) bool {
			return x < 3
		}, xs)

		asserts.Equal([]basics.Int{1, 2}, ToSlice(First(SUT)))
		asserts.Equal([]basics.Int{3, 4, 5}, ToSlice(Second(SUT)))
	})

	t.Run("Unzip", func(t *testing.T) {
		xs := FromSlice([]Tuple2[basics.Int, bool]{Pair(basics.Int(0), true), Pair(basics.Int(17), false), Pair(basics.Int(1337), true)})
		SUT := Unzip(xs)

		asserts.Equal([]basics.Int{0, 17, 1337}, ToSlice(First(SUT)))
		asserts.Equal([]bool{true, false, true}, ToSlice(Second(SUT)))
	})
}

func TestEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		SUT := Empty[basics.Int]()

		asserts.Equal(empty[basics.Int]{}, SUT)
	})
}
