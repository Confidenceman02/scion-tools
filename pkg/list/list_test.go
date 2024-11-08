package list

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	. "github.com/Confidenceman02/scion-tools/pkg/string"
	. "github.com/Confidenceman02/scion-tools/pkg/tuple"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmp(t *testing.T) {
	asserts := assert.New(t)

	t.Run(" Compare", func(t *testing.T) {
		x1 := fromSlice([]Int{Int(1), Int(2), Int(3)})
		x2 := fromSlice([]Int{Int(1), Int(2), Int(4)})

		asserts.Equal(LT{}, Compare(x1, x2))
	})

	t.Run("Compare: When empty", func(t *testing.T) {
		l1 := Empty[Int]()
		l2 := Empty[Int]()

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(EQ{}, Compare(l1, l2))
	})

	t.Run("Compare: When cons Int", func(t *testing.T) {
		l1 := Singleton[Int](1)
		l2 := Singleton[Int](1)
		l3 := Range(1, 10)
		l4 := Range(1, 10)
		l5 := Singleton[Int](2)
		l6 := Singleton[Int](1)
		l7 := Range(1, 11)
		l8 := Range(1, 10)
		l9 := Singleton[Int](1)
		l10 := Singleton[Int](2)
		l11 := Range(1, 10)
		l12 := Range(1, 11)
		l13 := Empty[Int]()
		l14 := Singleton[Int](2)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(0, l3.Cmp(l4))
		asserts.Equal(EQ{}, Compare(l1, l2))
		asserts.Equal(EQ{}, Compare(l3, l4))
		asserts.Equal(+1, l5.Cmp(l6))
		asserts.Equal(+1, l7.Cmp(l8))
		asserts.Equal(GT{}, Compare(l5, l6))
		asserts.Equal(GT{}, Compare(l7, l8))
		asserts.Equal(-1, l9.Cmp(l10))
		asserts.Equal(-1, l11.Cmp(l12))
		asserts.Equal(LT{}, Compare(l9, l10))
		asserts.Equal(LT{}, Compare(l11, l12))
		asserts.Equal(-1, l13.Cmp(l14))
		asserts.Equal(LT{}, Compare(l13, l14))
	})

	t.Run("Compare: When cons Float", func(t *testing.T) {
		l1 := Singleton[Float](1.0)
		l2 := Singleton[Float](1.0)
		l3 := Singleton[Float](1.1)
		l4 := Singleton[Float](1.0)
		l5 := fromSlice([]Float{1.1, 2.2})
		l6 := fromSlice([]Float{1.1, 2.1})
		l7 := Singleton[Float](1)
		l8 := Singleton[Float](2)
		l9 := fromSlice([]Float{1.0, 2.0})
		l10 := fromSlice([]Float{1.0, 2.1})
		l11 := Empty[Float]()
		l12 := Singleton[Float](2)

		asserts.Equal(0, l1.Cmp(l2))
		asserts.Equal(EQ{}, Compare(l1, l2))
		asserts.Equal(+1, l3.Cmp(l4))
		asserts.Equal(+1, l5.Cmp(l6))
		asserts.Equal(GT{}, Compare(l3, l4))
		asserts.Equal(GT{}, Compare(l5, l6))
		asserts.Equal(-1, l7.Cmp(l8))
		asserts.Equal(-1, l9.Cmp(l10))
		asserts.Equal(LT{}, Compare(l7, l8))
		asserts.Equal(LT{}, Compare(l9, l10))
		asserts.Equal(-1, l11.Cmp(l12))
		asserts.Equal(LT{}, Compare(l11, l12))
	})
}

func TestCreateFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[Int](10)

		asserts.Equal(&list[Int]{&_cons[Int]{a: 10, b: empty[Int]{}}}, SUT)
	})

	t.Run("Repeat", func(t *testing.T) {
		SUT := Repeat(2, 10)

		asserts.Equal(
			&list[int]{
				&_cons[int]{
					a: 10,
					b: &list[int]{
						&_cons[int]{
							a: 10,
							b: empty[int]{},
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
				&_cons[Int]{
					a: 2,
					b: &list[Int]{
						&_cons[Int]{
							a: 3,
							b: empty[Int]{},
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
				&_cons[int]{
					a: 20,
					b: ls,
				},
			},
			SUT,
		)
	})
}

func TestTransformFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map", func(t *testing.T) {
		SUT1 := fromSlice([]Float{Float(1), Float(4), Float(9)})
		SUT2 := fromSlice([]bool{true, false, true})

		asserts.Equal(&list[Float]{
			&_cons[Float]{
				a: 1,
				b: &list[Float]{
					&_cons[Float]{
						a: 2,
						b: &list[Float]{
							&_cons[Float]{
								a: 3,
								b: empty[Float]{},
							},
						},
					},
				},
			},
		},
			Map(Sqrt, SUT1),
		)

		asserts.Equal(&list[bool]{
			&_cons[bool]{
				a: false,
				b: &list[bool]{
					&_cons[bool]{
						a: true,
						b: &list[bool]{
							&_cons[bool]{
								a: false,
								b: empty[bool]{},
							},
						},
					},
				},
			},
		},
			Map(Not, SUT2),
		)
	})

	t.Run("IndexedMap", func(t *testing.T) {
		x1 := Range(0, 2)

		SUT := IndexedMap(Add, x1)

		asserts.Equal(
			&list[Int]{
				&_cons[Int]{
					a: Int(0),
					b: &list[Int]{
						&_cons[Int]{
							a: Int(2),
							b: &list[Int]{
								&_cons[Int]{
									a: Int(4),
									b: empty[Int]{},
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
			SUT := Foldl(Add, 0, ls)

			asserts.Equal(Int(6), SUT)
		})
		t.Run("Concat", func(t *testing.T) {
			ls := Range(1, 3)
			SUT := Foldl[Int, List[Int]](Cons, Empty[Int](), ls)

			asserts.Equal(
				&list[Int]{
					&_cons[Int]{
						a: 3,
						b: &list[Int]{
							&_cons[Int]{
								a: 2,
								b: &list[Int]{
									&_cons[Int]{
										a: 1,
										b: empty[Int]{},
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
				&_cons[Int]{
					a: 1,
					b: &list[Int]{
						&_cons[Int]{
							a: 2,
							b: &list[Int]{
								&_cons[Int]{
									a: 3,
									b: empty[Int]{},
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
		isEven := func(i Int) bool { return ModBy(2, i) == 0 }

		SUT := Filter(isEven, xs)

		asserts.Equal(
			&list[Int]{
				&_cons[Int]{
					a: Int(2),
					b: &list[Int]{
						&_cons[Int]{
							a: Int(4),
							b: &list[Int]{
								&_cons[Int]{
									a: Int(6),
									b: empty[Int]{},
								},
							},
						},
					},
				},
			},
			SUT,
		)
	})
	t.Run("FilterMap", func(t *testing.T) {
		xs := fromSlice([]String{"3", "hi", "12", "4th", "May"})
		SUT := FilterMap(ToInt, xs)

		asserts.Equal(
			&list[Int]{
				&_cons[Int]{
					a: Int(3),
					b: &list[Int]{
						&_cons[Int]{
							a: Int(12),
							b: empty[Int]{},
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
				&_cons[Int]{
					a: 3,
					b: &list[Int]{
						&_cons[Int]{
							a: 2,
							b: &list[Int]{
								&_cons[Int]{
									a: 1,
									b: empty[Int]{},
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
				&_cons[List[Int]]{
					a: &list[Int]{
						&_cons[Int]{
							a: 12,
							b: empty[Int]{},
						},
					},
				}}

			needle := &list[Int]{&_cons[Int]{a: 12, b: empty[Int]{}}}

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
	t.Run("Maximum", func(t *testing.T) {
		xs := fromSlice([]Comparable[Int]{Int(1), Int(4), Int(2)})
		SUT := Maximum(xs)

		asserts.Equal(maybe.Just[Int]{Value: Int(4)}, SUT)
	})
	t.Run("Maximum_UNSAFE", func(t *testing.T) {
		xs := fromSlice([]Int{Int(1), Int(2), Int(4)})
		SUT := Maximum_UNSAFE(xs)

		asserts.Equal(maybe.Just[Int]{Value: Int(4)}, SUT)
	})
	t.Run("Minimum", func(t *testing.T) {
		xs1 := fromSlice([]Comparable[Int]{Int(3), Int(2), Int(1)})
		SUT1 := Minimum(xs1)
		SUT2 := Minimum(Empty[Comparable[Int]]())

		asserts.Equal(maybe.Just[Int]{Value: Int(1)}, SUT1)
		asserts.Equal(maybe.Nothing{}, SUT2)
	})
	t.Run("Sum", func(t *testing.T) {
		xs1 := fromSlice([]Int{Int(1), Int(2), Int(3)})
		xs2 := fromSlice([]Int{Int(1), Int(1), Int(1)})
		SUT1 := Sum(xs1)
		SUT2 := Sum(xs2)
		SUT3 := Sum(Empty[Int]())

		asserts.Equal(Int(6), SUT1)
		asserts.Equal(Int(3), SUT2)
		asserts.Equal(Int(0), SUT3)
	})
	t.Run("Product", func(t *testing.T) {
		xs1 := fromSlice([]Int{Int(2), Int(2), Int(2)})
		xs2 := fromSlice([]Int{Int(3), Int(3), Int(3)})
		SUT1 := Product(xs1)
		SUT2 := Product(xs2)
		SUT3 := Product(Empty[Int]())

		asserts.Equal(Int(8), SUT1)
		asserts.Equal(Int(27), SUT2)
		asserts.Equal(Int(1), SUT3)

	})
}

func TestCombineFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Append with empty", func(t *testing.T) {
		xs := Empty[Int]()
		ys := Range(1, 3)
		SUT := Append(xs, ys)

		asserts.Equal(ys, SUT)
	})

	t.Run("Append", func(t *testing.T) {
		xs1 := Singleton[Int](1)
		ys1 := Range(2, 3)
		xs2 := fromSlice([]Int{Int(1), Int(1), Int(2)})
		ys2 := fromSlice([]Int{Int(3), Int(5), Int(8)})
		SUT1 := Tail(Append(xs1, ys1))
		SUT2 := Append(xs2, ys2)

		asserts.Equal(maybe.Just[List[Int]]{Value: ys1}, SUT1)
		asserts.Equal(&list[Int]{&_cons[Int]{a: 1, b: empty[Int]{}}}, xs1)
		asserts.Equal([]Int{1, 1, 2, 3, 5, 8}, toSlice(SUT2))
		// Structure sharing
		asserts.Equal(&list[Int]{
			&_cons[Int]{
				a: 1,
				b: &list[Int]{
					&_cons[Int]{
						a: 1,
						b: &list[Int]{
							&_cons[Int]{a: 2, b: ys2},
						},
					},
				},
			},
		},
			SUT2,
		)
	})

	t.Run("Concat", func(t *testing.T) {
		lists := fromSlice([]List[Int]{Range(1, 2), Range(3, 3), Range(4, 5)})
		SUT := Concat(lists)

		asserts.Equal([]Int{1, 2, 3, 4, 5}, toSlice(SUT))
	})

	t.Run("ConcatMap", func(t *testing.T) {
		fun := func(a Int) List[Int] { return Singleton(a + 1) }
		xs := Range(1, 3)
		SUT := ConcatMap(fun, xs)

		asserts.Equal([]Int{2, 3, 4}, toSlice(SUT))
		asserts.Equal([]Int{1, 2, 3}, toSlice(xs))
	})

	t.Run("Intersperse", func(t *testing.T) {
		xs := fromSlice([]string{"turtles", "turtles", "turtles"})
		SUT := Intersperse("on", xs)

		asserts.Equal([]string{"turtles", "on", "turtles", "on", "turtles"}, toSlice(SUT))
	})

	t.Run("Map2", func(t *testing.T) {
		t.Run("Lists with one member", func(t *testing.T) {
			a := Cons(1, Empty[Int]())
			b := Cons(1, Empty[Int]())
			SUT := Map2(Add, a, b)

			asserts.Equal(
				&list[Int]{
					&_cons[Int]{
						a: 2, b: Empty[Int]()}},
				SUT,
			)
		})
		t.Run("When one list empty", func(t *testing.T) {
			a := Range(1, 3)
			b := Empty[Int]()
			SUT := Map2(Add, a, b)

			asserts.Equal(empty[Int]{}, SUT)
		})
	})

	t.Run("Map3", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 3)
		xc := Range(1, 3)
		add3 := func(a Int, b Int, c Int) Int { return Add(a, Add(b, c)) }

		SUT := Map3(add3, xa, xb, xc)

		asserts.Equal(
			[]Int{3, 6, 9},
			toSlice(SUT),
		)
	})
	t.Run("Map3 with asymmetrical lists", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 2)
		xc := Range(1, 3)
		add3 := func(a Int, b Int, c Int) Int { return Add(a, Add(b, c)) }

		SUT := Map3(add3, xa, xb, xc)

		asserts.Equal(
			[]Int{3, 6},
			toSlice(SUT),
		)
	})

	t.Run("Map4", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 3)
		xc := Range(1, 3)
		xd := Range(1, 3)
		add4 := func(a Int, b Int, c Int, d Int) Int { return Add(a, Add(b, (Add(c, d)))) }

		SUT := Map4(add4, xa, xb, xc, xd)

		asserts.Equal(
			[]Int{4, 8, 12},
			toSlice(SUT),
		)
	})

	t.Run("Map5", func(t *testing.T) {
		xa := Range(1, 3)
		xb := Range(1, 3)
		xc := Range(1, 3)
		xd := Range(1, 3)
		xe := Range(1, 3)
		add5 := func(a Int, b Int, c Int, d Int, e Int) Int { return Add(a, Add(b, (Add(c, Add(d, e))))) }

		SUT := Map5(add5, xa, xb, xc, xd, xe)

		asserts.Equal(
			[]Int{5, 10, 15},
			toSlice(SUT),
		)
	})
}

func testSortFUnctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Sort", func(t *testing.T) {
		xs := fromSlice([]Comparable[String]{String("chuck"), String("alice"), String("bob")})
		SUT := Sort(xs)

		asserts.Equal([]String{"alice", "bob", "chuck"}, toSlice(SUT))
	})

	t.Run("Sort_UNSASFE", func(t *testing.T) {
		xs := fromSlice([]String{"chuck", "alice", "bob"})
		SUT := Sort_UNSAFE(xs)

		asserts.Equal([]String{"alice", "bob", "chuck"}, toSlice(SUT))
	})

	t.Run("SortBy", func(t *testing.T) {
		xs := fromSlice([]String{"chuck", "alice", "bob"})
		SUT := SortBy(func(s String) Comparable[String] { return s }, xs)

		asserts.Equal([]String{"alice", "bob", "chuck"}, toSlice(SUT))
	})

	t.Run("SortWith", func(t *testing.T) {
		xs := Range(1, 5)
		flippedComparison := func(a, b Int) Order {
			switch Compare(a, b) {
			case LT{}:
				return GT{}
			case EQ{}:
				return EQ{}
			default:
				return GT{}
			}
		}
		SUT := SortWith(flippedComparison, xs)

		asserts.Equal([]Int{5, 4, 3, 2, 1}, toSlice(SUT))
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
						&_cons[int]{
							a: 23,
							b: empty[int]{},
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

		asserts.Equal([]Int{1}, toSlice(SUT1))
		asserts.Equal([]Int{1, 2}, toSlice(SUT2))
		asserts.Equal([]Int{1, 2, 3}, toSlice(SUT3))
		asserts.Equal([]Int{1, 2, 3, 4}, toSlice(SUT4))
		asserts.Equal([]Int{1, 2, 3, 4, 5}, toSlice(SUT5))
		asserts.Equal([]Int{1}, toSlice(SUT6))
		asserts.Equal([]Int{1, 2}, toSlice(SUT7))
		asserts.Equal([]Int{1, 2, 3}, toSlice(SUT8))
		asserts.Equal([]Int{1, 2, 3, 4}, toSlice(SUT9))
		asserts.Equal([]Int{1, 2, 3, 4, 5}, toSlice(SUT10))
	})

	t.Run("Drop", func(t *testing.T) {
		xs := Range(1, 4)
		SUT := Drop(2, xs)

		asserts.Equal([]Int{3, 4}, toSlice(SUT))
	})

	t.Run("Partition", func(t *testing.T) {
		xs := Range(1, 5)
		SUT := Partition(func(x Int) bool {
			return x < 3
		}, xs)

		asserts.Equal([]Int{1, 2}, toSlice(First(SUT)))
		asserts.Equal([]Int{3, 4, 5}, toSlice(Second(SUT)))
	})

	t.Run("Unzip", func(t *testing.T) {
		xs := fromSlice([]Tuple2[Int, bool]{Pair(Int(0), true), Pair(Int(17), false), Pair(Int(1337), true)})
		SUT := Unzip(xs)

		asserts.Equal([]Int{0, 17, 1337}, toSlice(First(SUT)))
		asserts.Equal([]bool{true, false, true}, toSlice(Second(SUT)))
	})
}

func TestEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		SUT := Empty[Int]()

		asserts.Equal(empty[Int]{}, SUT)
	})
}
