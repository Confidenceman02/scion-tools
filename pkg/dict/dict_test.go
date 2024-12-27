package dict

import (
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	s "github.com/Confidenceman02/scion-tools/pkg/string"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(FromList(list.Empty[tuple.Tuple2[Int, any]]()), Empty[Int, any]())
	})
	t.Run("Singleton", func(t *testing.T) {
		expected := FromList(list.FromSlice(
			[]tuple.Tuple2[s.String, s.String]{
				tuple.Pair(s.String(s.String("k")), s.String("v")),
			}))
		SUT := Singleton(s.String("k"), s.String("v"))

		asserts.Equal(expected, SUT)
	})
	t.Run("Insert", func(t *testing.T) {
		expected := FromList(list.FromSlice(
			[]tuple.Tuple2[s.String, s.String]{
				tuple.Pair(s.String(s.String("k")), s.String("v")),
			}))
		SUT := Insert(s.String("k"), s.String("v"), Empty[s.String, s.String]())

		asserts.Equal(expected, SUT)
	})
	t.Run("Insert replace", func(t *testing.T) {
		expected := FromList(list.FromSlice(
			[]tuple.Tuple2[s.String, s.String]{
				tuple.Pair(s.String(s.String("k")), s.String("vv")),
			}))
		SUT := Insert(s.String("k"), s.String("vv"), Singleton(s.String("k"), s.String("v")))

		asserts.Equal(expected, SUT)
	})
	t.Run("Insert on Singleton", func(t *testing.T) {
		d := Singleton(Int(1), 2)
		d1 := Insert(Int(2), 2, d)

		asserts.Equal(Int(1), d.rbt().root.key)
		asserts.Equal(Int(1), d1.rbt().root.key)
		asserts.Nil(d.rbt().root.right)
		asserts.NotNil(d1.rbt().root.right)
		asserts.Equal(Int(2), d1.rbt().root.right.key)
	})
	t.Run("Insert on Empty", func(t *testing.T) {
		d := Empty[Int, Int]()
		d1 := Insert(Int(2), 2, d)

		asserts.Equal(&dict[Int, Int]{root: nil}, d.rbt())
		asserts.Equal(Int(2), d1.rbt().root.key)
	})
	t.Run("Insert into existing entry", func(t *testing.T) {
		d := Singleton(Int(10), 233)
		d1 := Insert(Int(10), 233, d)

		SUT := d1.rbt()

		asserts.Equal(233, d.rbt().root.value)
		asserts.Equal(233, SUT.root.value)
		asserts.NotSame(d.rbt().root, SUT.root)
	})
	t.Run("Insert with color pushdown", func(t *testing.T) {
		d := Singleton(Int(40), 1)
		d1 := Insert(Int(50), 2, d)
		d2 := Insert(Int(30), 3, d1)
		// Will cause color pushdown of parent node (40)
		d3 := Insert(Int(35), 3, d2)

		asserts.Equal(d1.rbt().root.right, d2.rbt().root.right)
		asserts.NotEqual(d2.rbt().root.right, d3.rbt().root.right)
	})
	t.Run("Update", func(t *testing.T) {
		expected := FromList(list.FromSlice(
			[]tuple.Tuple2[s.String, s.String]{
				tuple.Pair(s.String(s.String("k")), s.String("vv")),
			}))
		SUT := Update(s.String("k"), func(m maybe.Maybe[s.String]) maybe.Maybe[s.String] {
			return maybe.Just[s.String]{Value: s.String("vv")}
		}, Singleton(s.String("k"), s.String("v")))

		asserts.Equal(expected, SUT)
	})
	t.Run("Remove", func(t *testing.T) {
		expected := Empty[s.String, s.String]()
		SUT := Remove(s.String("k"), Singleton(s.String("k"), s.String("v")))

		asserts.Equal(expected, SUT)
	})
	t.Run("Remove not found", func(t *testing.T) {
		expected := Singleton(s.String("k"), s.String("v"))
		SUT := Remove(s.String("kk"), Singleton(s.String("k"), s.String("v")))

		asserts.Equal(expected, SUT)
	})
}

func TestMember(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Member on Singleton", func(t *testing.T) {
		d := Singleton(Int(10), 23)

		asserts.Equal(true, Member(Int(10), d))
		asserts.Equal(false, Member(Int(2), d))
	})

	t.Run("Member on Empty", func(t *testing.T) {
		d := Empty[Int, Int]()

		asserts.Equal(false, Member(Int(22), d))
		asserts.Equal(false, Member(Int(2), d))
	})
}

func TestIsEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty dict", func(t *testing.T) {
		d := Empty[Int, int]()
		SUT := IsEmpty(d)

		asserts.True(SUT)
	})

	t.Run("Singleton dict", func(t *testing.T) {
		d := Singleton(Int(100), 1)
		SUT := IsEmpty(d)

		asserts.False(SUT)
	})
}

func TestGet(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Get existing node", func(t *testing.T) {
		d := Singleton(Int(10), 23)
		SUT := Get(Int(10), d)

		asserts.Equal(maybe.Just[int]{Value: 23}, SUT)
	})

	t.Run("Get non-existing entry", func(t *testing.T) {
		d := Empty[Int, Int]()
		SUT := Get(Int(10), d)

		asserts.Equal(maybe.Nothing{}, SUT)
	})
}

func TestSize(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Size", func(t *testing.T) {
		xs := list.FromSlice([]Int{1, 2, 3, 4, 5})
		SUT := list.Foldl(
			func(i Int, d Dict[Int, Int]) Dict[Int, Int] { return Insert(i, i, d) },
			Empty[Int, Int](),
			xs,
		)

		asserts.Equal(Int(5), Size(SUT))
	})
}

func TestInsert(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(&dict[Int, struct{}]{root: nil}, Empty[Int, struct{}]())
	})

	t.Run("Insert on Empty has properties", func(t *testing.T) {
		d := Empty[Int, Int]()
		d1 := Insert(Int(1), 233, d)

		SUT := d1.rbt()

		asserts.Equal(&dict[Int, Int]{root: &node[Int, Int]{
			key:   Int(1),
			value: 233,
			color: black,
			left:  nil,
			right: nil}},
			SUT,
		)
	})

	t.Run("Insert on Singleton right side", func(t *testing.T) {
		d := Singleton[Int, Int](Int(1), 1)
		d1 := Insert(Int(2), 2, d)

		SUT := d1.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(red, SUT.root.right.color)
	})

	t.Run("Insert into existing entry", func(t *testing.T) {
		d := Singleton(Int(10), 233)
		d1 := Insert(Int(10), 100, d)

		SUT := d1

		asserts.Equal(&dict[Int, int]{root: &node[Int, int]{
			key:   Int(10),
			value: 100,
			color: black,
			left:  nil,
			right: nil},
		}, SUT)
	})

	t.Run("LL Single right rotation", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		Insert(Int(30), 3, d1)

		asserts.Nil(d1.rbt().root.right)
		asserts.Equal(Int(40), d1.rbt().root.left.key)
	})

	t.Run("LR -> RR rotation", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(45), 3, d1)

		// d1
		asserts.Equal(Int(50), d1.rbt().root.key)
		asserts.Nil(d1.rbt().root.right)
		asserts.Equal(Int(40), d1.rbt().root.left.key)
		asserts.Nil(d1.rbt().root.left.left)
		asserts.Nil(d1.rbt().root.left.right)

		// d2
		asserts.Equal(Int(45), d2.rbt().root.key)
		asserts.Equal(Int(50), d2.rbt().root.right.key)
		asserts.Nil(d2.rbt().root.right.right)
		asserts.Equal(Int(40), d2.rbt().root.left.key)
		asserts.Nil(d2.rbt().root.left.left)
	})

	t.Run("RR Single left rotation", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(60), 2, d)
		Insert(Int(70), 3, d1)

		asserts.Nil(d1.rbt().root.left)
		asserts.Equal(Int(60), d1.rbt().root.right.key)
	})

	t.Run("LR -> LL rotation", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(60), 2, d)
		d2 := Insert(Int(55), 3, d1)

		// d1
		asserts.Equal(Int(50), d1.rbt().root.key)
		asserts.Nil(d1.rbt().root.left)
		asserts.Equal(Int(60), d1.rbt().root.right.key)
		asserts.Nil(d1.rbt().root.right.right)
		asserts.Nil(d1.rbt().root.right.left)

		// d2
		asserts.Equal(Int(55), d2.rbt().root.key)
		asserts.Equal(Int(60), d2.rbt().root.right.key)
		asserts.Nil(d2.rbt().root.left.left)
		asserts.Equal(Int(50), d2.rbt().root.left.key)
		asserts.Nil(d2.rbt().root.right.right)
	})

	t.Run("granparent color pushdown", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(45), 3, d1)
		d3 := Insert(Int(30), 3, d2)

		// d2
		asserts.Equal(Int(50), d2.rbt().root.right.key)
		asserts.Equal(red, d2.rbt().root.right.color)
		asserts.Nil(d2.rbt().root.left.left)

		// d3
		asserts.Equal(Int(50), d3.rbt().root.right.key)
		asserts.Equal(black, d3.rbt().root.right.color)
		asserts.Equal(black, d3.rbt().root.left.color)
		asserts.Equal(red, d3.rbt().root.left.left.color)
	})

	t.Run("LL Single right rotation", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(30), 3, d1)

		SUT := d2.rbt()

		asserts.Equal(Int(40), SUT.root.key)
		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(50), SUT.root.right.key)
		asserts.Equal(red, SUT.root.right.color)
		asserts.Equal(Int(30), SUT.root.left.key)
		asserts.Equal(red, SUT.root.left.color)
	})

	t.Run("RR Single right rotation", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(60), 2, d)
		d2 := Insert(Int(70), 3, d1)

		SUT := d2.rbt()

		asserts.Equal(Int(60), SUT.root.key)
		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(50), SUT.root.left.key)
		asserts.Equal(red, SUT.root.left.color)
		asserts.Equal(Int(70), SUT.root.right.key)
		asserts.Equal(red, SUT.root.right.color)
	})

	t.Run("LR double red, red uncle", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		// Left
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(60), 3, d1)
		d3 := Insert(Int(45), 4, d2)

		SUT := d3.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(black, SUT.root.left.color)
		asserts.Equal(black, SUT.root.right.color)
		asserts.Equal(red, SUT.root.left.right.color)
	})

	t.Run("LR double red, black uncle", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		// Left
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(45), 3, d1)

		SUT := d2.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(45), SUT.root.key)
		asserts.Equal(red, SUT.root.left.color)
		asserts.Equal(Int(40), SUT.root.left.key)
		asserts.Equal(red, SUT.root.right.color)
		asserts.Equal(Int(50), SUT.root.right.key)
	})

	t.Run("RL double red, red uncle", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(60), 2, d)
		d2 := Insert(Int(40), 3, d1)
		d3 := Insert(Int(55), 4, d2)

		SUT := d3.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(black, SUT.root.left.color)
		asserts.Equal(black, SUT.root.right.color)
		asserts.Equal(red, SUT.root.right.left.color)
	})

	t.Run("RL double red, black uncle", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(60), 2, d)
		d2 := Insert(Int(55), 4, d1)

		SUT := d2.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(55), SUT.root.key)
		asserts.Equal(red, SUT.root.left.color)
		asserts.Equal(Int(50), SUT.root.left.key)
		asserts.Equal(red, SUT.root.right.color)
		asserts.Equal(Int(60), SUT.root.right.key)
	})

	t.Run("test the following inserts 7,5,10,20,15", func(t *testing.T) {
		d := Singleton(Int(7), 1)
		d1 := Insert(Int(5), 2, d)
		d2 := Insert(Int(10), 3, d1)
		d3 := Insert(Int(20), 3, d2)
		d4 := Insert(Int(15), 3, d3)

		SUT := d4.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(7), SUT.root.key)
		asserts.Equal(black, SUT.root.right.color)
		asserts.Equal(Int(15), SUT.root.right.key)
		asserts.Equal(red, SUT.root.right.right.color)
		asserts.Equal(Int(20), SUT.root.right.right.key)
		asserts.Equal(red, SUT.root.right.left.color)
		asserts.Equal(Int(10), SUT.root.right.left.key)
	})

	t.Run("test the following inserts 10,15,5,0,2", func(t *testing.T) {
		d := Singleton(Int(10), 1)
		d1 := Insert(Int(15), 2, d)
		d2 := Insert(Int(5), 3, d1)
		d3 := Insert(Int(0), 3, d2)
		d4 := Insert(Int(2), 3, d3)

		SUT := d4.rbt()

		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(10), SUT.root.key)
		asserts.Equal(black, SUT.root.left.color)
		asserts.Equal(Int(2), SUT.root.left.key)
		asserts.Equal(red, SUT.root.left.left.color)
		asserts.Equal(Int(0), SUT.root.left.left.key)
		asserts.Equal(red, SUT.root.left.right.color)
		asserts.Equal(Int(5), SUT.root.left.right.key)
	})

	t.Run("Structure sharing right subtree", func(t *testing.T) {
		d := Singleton(Int(40), 1)
		d1 := Insert(Int(50), 2, d)
		d2 := Insert(Int(30), 3, d1)
		d3 := Insert(Int(35), 3, d1)

		asserts.Equal(d1.rbt().root.right, d2.rbt().root.right)
		asserts.Equal(d2.rbt().root.right, d3.rbt().root.right)
	})
}

func TestRemove(t *testing.T) {
	asserts := assert.New(t)

	t.Run("remove Empty Dict", func(t *testing.T) {
		d := Empty[Int, Int]()
		d1 := Remove(Int(50), d)

		asserts.Equal(&d, &d1)
	})

	t.Run("remove Singleton key that doesn't exist", func(t *testing.T) {
		d := Singleton[Int, Int](Int(1), 1)
		d1 := Remove(Int(50), d)

		// Pointers match
		asserts.Equal(&d, &d1)
	})

	t.Run("remove Singleton", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Remove(Int(50), d)

		asserts.NotNil(d.rbt().root)
		asserts.Equal(Int(50), d.rbt().root.key)
		asserts.Nil(d1.rbt().root)
	})

	t.Run("remove childless red leaf", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(60), 3, d1)
		d3 := Remove(Int(40), d2)

		asserts.NotNil(d1.rbt().root.left)
		asserts.NotNil(d2.rbt().root.left)
		asserts.Nil(d3.rbt().root.left)

		// Structure Sharing
		asserts.True(d2.rbt().root.right == d3.rbt().root.right)

		// Different root d2, d3
		asserts.NotEqual(d.rbt().root, d1.rbt().root)
		asserts.NotEqual(d1.rbt().root, d2.rbt().root)
		asserts.NotEqual(d2.rbt().root, d3.rbt().root)
	})

	t.Run("black leaf, LEFT | BLACK sibling, red near nephew, BLACK distant nephew", func(t *testing.T) {
		var tree Dict[Int, Comparable[Int]]
		tree = &dict[Int, Comparable[Int]]{
			root: &node[Int, Comparable[Int]]{
				key:   Int(40),
				value: Int(1),
				color: black,
				left:  &node[Int, Comparable[Int]]{key: Int(30), value: Int(2), color: black, left: nil, right: nil},
				right: &node[Int, Comparable[Int]]{
					key:   Int(50),
					value: Int(3),
					color: black,
					left:  &node[Int, Comparable[Int]]{key: Int(45), value: Int(4), color: red, left: nil, right: nil},
					right: nil,
				},
			},
		}

		SUT := Remove(Int(30), tree).rbt()

		asserts.Equal(Int(40), tree.rbt().root.key)
		asserts.Equal(Int(45), SUT.root.key)
		asserts.Equal(Int(50), SUT.root.right.key)
		asserts.Equal(Int(40), SUT.root.left.key)
		asserts.NotEqual(tree.rbt().root.right, SUT.root.right)
	})

	t.Run("black leaf, RIGHT | BLACK sibling, red near nephew, BLACK distant nephew", func(t *testing.T) {
		var tree Dict[Int, Comparable[Int]]
		tree = &dict[Int, Comparable[Int]]{
			root: &node[Int, Comparable[Int]]{
				key:   Int(40),
				value: Int(1),
				color: black,
				left: &node[Int, Comparable[Int]]{
					key:   Int(30),
					value: Int(2),
					color: black,
					left:  nil,
					right: &node[Int, Comparable[Int]]{key: Int(35), value: Int(4), color: red, left: nil, right: nil},
				},
				right: &node[Int, Comparable[Int]]{
					key:   Int(50),
					value: Int(3),
					color: black,
					left:  nil,
					right: nil,
				},
			},
		}

		SUT := Remove(Int(50), tree).rbt()

		asserts.Equal(Int(40), tree.rbt().root.key)
		asserts.Equal(Int(35), SUT.root.key)
		asserts.Equal(Int(40), SUT.root.right.key)
		asserts.Equal(Int(30), SUT.root.left.key)
		asserts.NotEqual(tree.rbt().root.left, SUT.root.left)
	})

	t.Run("black leaf, RIGHT | red sibling | BLACK near nephew | BLACK distant nephew", func(t *testing.T) {
		var tree Dict[Int, Int]
		tree = &dict[Int, Int]{
			root: &node[Int, Int]{
				key:   Int(50),
				value: 1,
				color: black,
				left: &node[Int, Int]{
					key:   Int(40),
					value: 2,
					color: red,
					left:  &node[Int, Int]{key: Int(35), value: 5, color: black, left: nil, right: nil},
					right: &node[Int, Int]{
						key:   Int(45),
						value: 6,
						color: black,
						left:  nil,
						right: nil,
					},
				},
				right: &node[Int, Int]{
					key:   Int(60),
					value: 3,
					color: black,
					left:  nil,
					right: nil,
				},
			},
		}

		SUT := Remove(Int(60), tree).rbt()

		// Removes node
		asserts.Nil(SUT.root.right.right)

		asserts.Equal(Int(50), tree.rbt().root.key)
		asserts.Equal(Int(40), SUT.root.key)
		asserts.Equal(black, SUT.root.color)
		asserts.Equal(Int(50), SUT.root.right.key)
		asserts.Equal(Int(45), SUT.root.right.left.key)
		asserts.Equal(red, SUT.root.right.left.color)

		// Structure sharing
		asserts.True(tree.rbt().root.left.left == SUT.root.left)
	})

	t.Run("black node, LEFT | red child, LEFT | NIL child, RIGHT", func(t *testing.T) {
		var tree Dict[Int, Int]
		tree = &dict[Int, Int]{
			root: &node[Int, Int]{
				key:   Int(50),
				value: 1,
				color: black,
				left: &node[Int, Int]{
					key:   Int(40),
					color: black,
					value: 3,
					left:  nil,
					right: &node[Int, Int]{key: Int(45), color: red, value: 6, left: nil, right: nil}},
				right: &node[Int, Int]{key: Int(60), color: black, value: 2, left: nil, right: nil},
			},
		}

		SUT := Remove(Int(40), tree).rbt()

		asserts.Nil(SUT.root.left.right)
		asserts.Equal(Int(45), SUT.root.left.key)

		// Structure Sharing
		asserts.True(tree.rbt().root.right == SUT.root.right)
	})

	t.Run("black node, RIGHT | red child, LEFT | NIL child, RIGHT", func(t *testing.T) {
		var tree Dict[Int, Int]
		tree = &dict[Int, Int]{
			root: &node[Int, Int]{
				key:   Int(50),
				value: 1,
				color: black,
				right: &node[Int, Int]{
					key:   Int(60),
					color: black,
					value: 3,
					left:  &node[Int, Int]{key: Int(55), color: red, value: 6, left: nil, right: nil},
					right: nil,
				},
				left: &node[Int, Int]{key: Int(40), color: black, value: 2, left: nil, right: nil},
			},
		}

		SUT := Remove(Int(60), tree).rbt()

		asserts.Nil(SUT.root.right.left)
		asserts.Equal(Int(55), SUT.root.right.key)

		// Structure Sharing
		asserts.True(tree.rbt().root.left == SUT.root.left)
	})

	t.Run("Removes root node with 2 red children", func(t *testing.T) {
		var tree Dict[Int, Int]
		tree = &dict[Int, Int]{
			root: &node[Int, Int]{
				key:   Int(50),
				color: black,
				value: 1,
				left:  &node[Int, Int]{key: Int(40), color: red, value: 2, left: nil, right: nil},
				right: &node[Int, Int]{key: Int(60), color: red, value: 3, left: nil, right: nil},
			},
		}

		SUT := Remove(Int(50), tree).rbt()

		asserts.Equal(Int(60), SUT.root.key)
		asserts.Nil(SUT.root.right)
		asserts.Equal(Int(40), SUT.root.left.key)

		// Structure Sharing
		asserts.Equal(tree.rbt().root.left, SUT.root.left)
		asserts.NotEqual(tree.rbt().root, SUT.root)
	})

	t.Run("Removes red right leaf node with no children", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(60), 3, d1)
		d3 := Remove(Int(60), d2)

		SUT := d3.rbt()

		asserts.Nil(SUT.root.right)

		// Structure Sharing
		asserts.True(d1.rbt().root.left == SUT.root.left)
		asserts.True(d1.rbt().root != SUT.root)
	})

	t.Run("Removes a red left node with no children", func(t *testing.T) {
		d := Singleton(Int(50), 1)
		d1 := Insert(Int(40), 2, d)
		d2 := Insert(Int(60), 3, d1)
		d3 := Remove(Int(40), d2)

		SUT := d3.rbt()

		asserts.Nil(SUT.root.left)
	})

	t.Run("Solve rbt | remove 50,20,100,90,40,60,70,10,30,80", func(t *testing.T) {
		// https://www.youtube.com/watch?v=PgO_Xj7DC1A&t=16s
		var tree Dict[Int, Int]
		tree = &dict[Int, Int]{
			root: &node[Int, Int]{
				key:   Int(40),
				value: 1,
				color: black,
				left: &node[Int, Int]{
					key:   Int(20),
					value: 2,
					color: black,
					left:  &node[Int, Int]{key: Int(10), value: 3, color: black, left: nil, right: nil},
					right: &node[Int, Int]{key: Int(30), value: 4, color: black, left: nil, right: nil},
				},
				right: &node[Int, Int]{
					key:   Int(60),
					value: 5,
					color: black,
					left:  &node[Int, Int]{key: Int(50), value: 6, color: black, left: nil, right: nil},
					right: &node[Int, Int]{
						key:   Int(80),
						value: 7,
						color: red,
						left:  &node[Int, Int]{key: Int(70), value: 8, color: black, left: nil, right: nil},
						right: &node[Int, Int]{
							key:   Int(90),
							value: 9,
							color: black,
							left:  nil,
							right: &node[Int, Int]{key: Int(100), value: 10, color: red, left: nil, right: nil},
						},
					},
				},
			},
		}

		// REMOVE 50
		tree1 := Remove(Int(50), tree)
		SUT1 := tree1.rbt()

		asserts.Equal(Int(40), SUT1.root.key)
		asserts.Equal(black, SUT1.root.color)

		// Structure sharing
		// 90
		asserts.True(tree.rbt().root.right.right.right == SUT1.root.right.right)
		// 100
		asserts.True(tree.rbt().root.right.right.right.right == SUT1.root.right.right.right)
		// 60
		asserts.Equal(tree.rbt().root.right.key, SUT1.root.right.left.key)
		asserts.NotEqual(tree.rbt().root.right, SUT1.root.right.left)
		// 70
		asserts.Equal(tree.rbt().root.right.right.left.key, SUT1.root.right.left.right.key)
		asserts.NotEqual(tree.rbt().root.right.right.left, SUT1.root.right.left.right)

		// REMOVE 20
		tree2 := Remove(Int(20), tree1)
		SUT2 := tree2.rbt()

		asserts.Equal(Int(40), SUT2.root.key)
		asserts.Equal(black, SUT2.root.color)

		// Structure sharing
		// 60
		asserts.True(SUT1.root.right.left == SUT2.root.right.left)
		// 90
		asserts.True(SUT1.root.right.right == SUT2.root.right.right)
		// 70
		asserts.True(SUT1.root.right.left.right == SUT2.root.right.left.right)
		// 100
		asserts.True(SUT1.root.right.right.right == SUT2.root.right.right.right)

		// REMOVE 100
		tree3 := Remove(Int(100), tree2)
		SUT3 := tree3.rbt()

		// Different root
		asserts.True(SUT2.root != SUT3.root)

		// Struture sharing
		// 50
		asserts.True(SUT2.root.left == SUT3.root.left)
		// 10
		asserts.True(SUT2.root.left.left == SUT3.root.left.left)

		// REMOVE 90
		tree4 := Remove(Int(90), tree3)
		SUT4 := tree4.rbt()

		asserts.Equal(Int(40), SUT4.root.key)
		asserts.Equal(Int(70), SUT4.root.right.key)
		asserts.Equal(red, SUT4.root.right.color)
		asserts.Equal(Int(80), SUT4.root.right.right.key)
		asserts.Equal(black, SUT4.root.right.right.color)
		asserts.Equal(Int(60), SUT4.root.right.left.key)
		asserts.Equal(black, SUT4.root.right.left.color)
		asserts.Nil(SUT4.root.right.left.left)
		asserts.Nil(SUT4.root.right.left.right)

		// REMOVE 40
		tree5 := Remove(Int(40), tree4)
		SUT5 := tree5.rbt()

		asserts.Equal(Int(60), SUT5.root.key)
		asserts.Equal(black, SUT5.root.color)
		asserts.Equal(Int(70), SUT5.root.right.key)
		asserts.Equal(black, SUT5.root.right.color)
		asserts.Equal(Int(80), SUT5.root.right.right.key)
		asserts.Equal(red, SUT5.root.right.right.color)

		// Structure sharing
		asserts.True(tree5.rbt().root.left == SUT5.root.left)
		asserts.True(tree5.rbt().root.left.left == SUT5.root.left.left)

		// REMOVE 60
		tree6 := Remove(Int(60), tree5)
		SUT6 := tree6.rbt()

		asserts.Equal(Int(70), SUT6.root.key)
		asserts.Equal(black, SUT6.root.color)
		asserts.Equal(Int(80), SUT6.root.right.key)
		asserts.Equal(black, SUT6.root.right.color)

		// REMOVE 70
		tree7 := Remove(Int(70), tree6)
		SUT7 := tree7.rbt()

		asserts.Equal(Int(30), SUT7.root.key)
		asserts.Equal(black, SUT7.root.color)
		asserts.Equal(Int(80), SUT7.root.right.key)
		asserts.Equal(black, SUT7.root.right.color)

		// Structure sharing
		asserts.True(tree7.rbt().root.left == SUT7.root.left)

		// REMOVE 10
		tree8 := Remove(Int(10), tree7)
		SUT8 := tree8.rbt()

		asserts.Equal(Int(30), SUT8.root.key)
		asserts.Equal(black, SUT8.root.color)
		asserts.Equal(Int(80), SUT8.root.right.key)
		asserts.Equal(red, SUT8.root.right.color)
		asserts.Nil(SUT8.root.left)

		// REMOVE 30
		tree9 := Remove(Int(30), tree8)
		SUT9 := tree9.rbt()

		asserts.Equal(Int(80), SUT9.root.key)
		asserts.Equal(black, SUT9.root.color)
		asserts.Nil(SUT9.root.right)

		// REMOVE 80
		tree10 := Remove(Int(80), tree9)
		SUT10 := tree10.rbt()

		asserts.Nil(SUT10.root)
	})
}

type User struct {
	Name string
	Age  int
}

func TestListFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Keys", func(t *testing.T) {
		xs := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(Int(1)), 1),
			tuple.Pair(Int(Int(2)), 2),
			tuple.Pair(Int(Int(3)), 3),
		})
		dxs := FromList(xs)
		SUT := Keys(dxs)

		asserts.Equal([]Int{Int(1), Int(2), Int(3)}, list.ToSlice(SUT))
	})
	t.Run("Values", func(t *testing.T) {
		xs := list.FromSlice([]tuple.Tuple2[Int, string]{
			tuple.Pair(Int(Int(1)), "Alice"),
			tuple.Pair(Int(Int(2)), "Bob"),
		})
		dxs := FromList(xs)
		SUT := Values(dxs)

		asserts.Equal([]string{"Alice", "Bob"}, list.ToSlice(SUT))
	})

	t.Run("ToList", func(t *testing.T) {
		xs := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(2), 2),
		})
		SUT := ToList(FromList(xs))

		asserts.Equal([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(2), 2),
		},
			list.ToSlice(SUT),
		)
	})

	t.Run("FromList", func(t *testing.T) {
		xs := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(2), 2),
			tuple.Pair(Int(3), 3),
		})
		SUT := FromList(xs)

		asserts.Equal([]Int{Int(1), Int(2), Int(3)}, list.ToSlice(Keys(SUT)))
	})
}

func insertBoth[K Comparable[K], B any](key K, leftVal list.List[B], rightVal list.List[B], dict Dict[K, list.List[B]]) Dict[K, list.List[B]] {
	return Insert(key, list.Append(leftVal, rightVal), dict)
}
func TestTransformFunctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map", func(t *testing.T) {
		xs := list.FromSlice([]tuple.Tuple2[Int, int]{tuple.Pair(Int(1), 1), tuple.Pair(Int(2), 2), tuple.Pair(Int(3), 3)})
		d := FromList(xs)
		mapper := func(k Int, v int) int { return (v + 1) }
		SUT := Map(mapper, d)

		asserts.Equal([]int{2, 3, 4}, list.ToSlice(Values(SUT)))
	})

	t.Run("Foldl", func(t *testing.T) {
		xs := list.FromSlice([]tuple.Tuple2[Int, string]{
			tuple.Pair(Int(1), "hello"),
			tuple.Pair(Int(2), "my"),
			tuple.Pair(Int(3), "name"),
			tuple.Pair(Int(4), "is"),
		})
		d := FromList(xs)
		SUT := Foldl(func(k Int, v string, acc string) string {
			return acc + v + "pause"
		},
			"",
			d,
		)

		asserts.Equal("hellopausemypausenamepauseispause", SUT)
	})

	t.Run("Foldr", func(t *testing.T) {
		getAges := func(_ Int, user User, ages list.List[Int]) list.List[Int] {
			return list.Cons(Int(user.Age), ages)
		}
		user1 := User{
			Name: "Jdawg",
			Age:  42,
		}
		user2 := User{
			Name: "TDawg",
			Age:  30,
		}
		users := list.FromSlice([]User{user1, user2})
		d := list.Foldl(func(usr User, acc Dict[Int, User]) Dict[Int, User] {
			return Insert(Int(usr.Age), usr, acc)
		},
			Empty[Int, User](),
			users,
		)
		SUT := Foldr(getAges, list.Empty[Int](), d)

		asserts.Equal([]Int{30, 42}, list.ToSlice(SUT))
	})

	t.Run("Filter", func(t *testing.T) {
		xxs := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(2), 2),
			tuple.Pair(Int(3), 3),
		})
		d := FromList(xxs)
		filterer := func(k Int, v int) bool { return v > 2 }
		SUT := Filter(filterer, d)

		asserts.Equal([]int{3}, list.ToSlice(Values(SUT)))
	})

	t.Run("Partition", func(t *testing.T) {
		xxs := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(2), 2),
			tuple.Pair(Int(3), 3),
		})
		d := FromList(xxs)
		isGood := func(k Int, v int) bool { return v > 2 }
		SUT := Partition(isGood, d)

		asserts.Equal([]int{3}, list.ToSlice(Values(tuple.First(SUT))))
		asserts.Equal([]int{1, 2}, list.ToSlice(Values(tuple.Second(SUT))))
	})

	t.Run("Union", func(t *testing.T) {
		xs1 := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 2),
			tuple.Pair(Int(2), 2),
			tuple.Pair(Int(3), 3),
		})
		xs2 := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(2), 2),
			tuple.Pair(Int(3), 3),
		})
		d1 := FromList(xs1)
		d2 := FromList(xs2)

		SUT := Union(d1, d2)

		asserts.Equal([]int{2, 2, 3}, list.ToSlice(Values(SUT)))
	})

	t.Run("Intersect", func(t *testing.T) {
		xs1 := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 23),
			tuple.Pair(Int(2), 2),
			tuple.Pair(Int(3), 3),
		})
		xs2 := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(5), 2),
			tuple.Pair(Int(6), 3),
		})

		d1 := FromList(xs1)
		d2 := FromList(xs2)

		SUT := Intersect(d1, d2)

		asserts.Equal([]Int{Int(1)}, list.ToSlice(Keys(SUT)))
		asserts.Equal([]int{23}, list.ToSlice(Values(SUT)))
	})

	t.Run("Diff", func(t *testing.T) {
		xs1 := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(23), 23),
			tuple.Pair(Int(5), 2),
			tuple.Pair(Int(6), 3),
		})
		xs2 := list.FromSlice([]tuple.Tuple2[Int, int]{
			tuple.Pair(Int(1), 1),
			tuple.Pair(Int(5), 2),
			tuple.Pair(Int(6), 3),
		})

		d1 := FromList(xs1)
		d2 := FromList(xs2)

		SUT := Diff(d1, d2)

		asserts.Equal([]Int{Int(23)}, list.ToSlice(Keys(SUT)))
		asserts.Equal([]int{23}, list.ToSlice(Values(SUT)))
	})

	t.Run("Merge", func(t *testing.T) {
		s1 := Insert(s.String("u1"), list.Singleton(Int(1)), Empty[s.String, list.List[Int]]())
		s2 := Insert(s.String("u2"), list.Singleton(Int(2)), Empty[s.String, list.List[Int]]())
		s23 := Insert(s.String("u2"), list.Singleton(Int(3)), Empty[s.String, list.List[Int]]())
		b1 := FromList(list.Map(func(i Int) tuple.Tuple2[Int, list.List[Int]] {
			return tuple.Pair(i, list.Singleton(i))
		},
			list.Range(1, 10),
		))
		b2 := FromList(list.Map(func(i Int) tuple.Tuple2[Int, list.List[Int]] {
			return tuple.Pair(i, list.Singleton(i))
		},
			list.Range(5, 15),
		))

		t.Run("merge empties", func(t *testing.T) {
			SUT := Merge(
				Insert[s.String, list.List[any]],
				insertBoth[s.String, any],
				Insert[s.String, list.List[any]],
				Empty[s.String, list.List[any]](),
				Empty[s.String, list.List[any]](),
				Empty[s.String, list.List[any]](),
			)

			asserts.True(IsEmpty(SUT))
		})
		t.Run("merge singletons in order", func(t *testing.T) {
			SUT := Merge(
				Insert[s.String, list.List[Int]],
				insertBoth[s.String, Int],
				Insert[s.String, list.List[Int]],
				s1,
				s2,
				Empty[s.String, list.List[Int]](),
			)

			expected := list.FromSlice([]tuple.Tuple2[s.String, list.List[Int]]{
				tuple.Pair(s.String("u1"), list.Singleton(Int(1))),
				tuple.Pair(s.String("u2"), list.Singleton(Int(2))),
			})

			asserts.Equal(expected, ToList(SUT))
		})
		t.Run("merge singletons out of order", func(t *testing.T) {
			SUT := Merge(
				Insert[s.String, list.List[Int]],
				insertBoth[s.String, Int],
				Insert[s.String, list.List[Int]],
				s2,
				s1,
				Empty[s.String, list.List[Int]](),
			)

			expected := list.FromSlice([]tuple.Tuple2[s.String, list.List[Int]]{
				tuple.Pair(s.String("u1"), list.Singleton(Int(1))),
				tuple.Pair(s.String("u2"), list.Singleton(Int(2))),
			})

			asserts.Equal(expected, ToList(SUT))
		})
		t.Run("merge with duplicate key", func(t *testing.T) {
			SUT := Merge(
				Insert[s.String, list.List[Int]],
				insertBoth[s.String, Int],
				Insert[s.String, list.List[Int]],
				s2,
				s23,
				Empty[s.String, list.List[Int]](),
			)

			expected := list.FromSlice([]tuple.Tuple2[s.String, list.List[Int]]{
				tuple.Pair(s.String("u2"), list.FromSlice([]Int{2, 3})),
			})

			asserts.Equal(expected, ToList(SUT))
		})
		t.Run("partially overlapping", func(t *testing.T) {
			SUT := Merge(
				Insert[Int, list.List[Int]],
				insertBoth[Int, Int],
				Insert[Int, list.List[Int]],
				b1,
				b2,
				Empty[Int, list.List[Int]](),
			)

			expected := list.FromSlice([]tuple.Tuple2[Int, list.List[Int]]{
				tuple.Pair(Int(1), list.FromSlice([]Int{Int(1)})),
				tuple.Pair(Int(2), list.FromSlice([]Int{Int(2)})),
				tuple.Pair(Int(3), list.FromSlice([]Int{Int(3)})),
				tuple.Pair(Int(4), list.FromSlice([]Int{Int(4)})),
				tuple.Pair(Int(5), list.FromSlice([]Int{Int(5), Int(5)})),
				tuple.Pair(Int(6), list.FromSlice([]Int{Int(6), Int(6)})),
				tuple.Pair(Int(7), list.FromSlice([]Int{Int(7), Int(7)})),
				tuple.Pair(Int(8), list.FromSlice([]Int{Int(8), Int(8)})),
				tuple.Pair(Int(9), list.FromSlice([]Int{Int(9), Int(9)})),
				tuple.Pair(Int(10), list.FromSlice([]Int{Int(10), Int(10)})),
				tuple.Pair(Int(11), list.FromSlice([]Int{Int(11)})),
				tuple.Pair(Int(12), list.FromSlice([]Int{Int(12)})),
				tuple.Pair(Int(13), list.FromSlice([]Int{Int(13)})),
				tuple.Pair(Int(14), list.FromSlice([]Int{Int(14)})),
				tuple.Pair(Int(15), list.FromSlice([]Int{Int(15)})),
			})

			asserts.Equal(expected, ToList(SUT))
		})
	})
}
