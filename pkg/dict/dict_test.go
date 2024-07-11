package dict

import (
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(&dict[int, struct{}]{root: nil}, Empty[int, struct{}]())
	})

	t.Run("Singleton", func(t *testing.T) {
		d := Singleton[int, struct{}](1, struct{}{})
		SUT := d.rbt()
		asserts.Equal(&dict[int, struct{}]{
			root: &node[int, struct{}]{
				key:   1,
				value: struct{}{},
				color: BLACK,
				left:  nil,
				right: nil},
		},
			SUT,
		)
	})
}

func TestMember(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Member on Singleton", func(t *testing.T) {
		d := Singleton(10, 23)

		asserts.Equal(true, Member(10, d))
		asserts.Equal(false, Member(2, d))
	})

	t.Run("Member on Empty", func(t *testing.T) {
		d := Empty[int, int]()

		asserts.Equal(false, Member(22, d))
		asserts.Equal(false, Member(2, d))
	})
}

func TestIsEmpty(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty dict", func(t *testing.T) {
		d := Empty[int, int]()
		SUT := IsEmpty(d)

		asserts.True(SUT)
	})

	t.Run("Singleton dict", func(t *testing.T) {
		d := Singleton[int, int](100, 1)
		SUT := IsEmpty(d)

		asserts.False(SUT)
	})
}

func TestGet(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Get existing node", func(t *testing.T) {
		d := Singleton(10, 23)
		SUT := Get(10, d)

		asserts.Equal(maybe.Just[int]{Value: 23}, SUT)
	})

	t.Run("Get non-existing entry", func(t *testing.T) {
		d := Empty[int, int]()
		SUT := Get(10, d)

		asserts.Equal(maybe.Nothing{}, SUT)
	})
}

func TestInsert(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Insert on Singleton", func(t *testing.T) {
		d := Singleton(1, 2)
		d1 := Insert(2, 2, d)

		asserts.Equal(1, d.rbt().root.key)
		asserts.Equal(1, d1.rbt().root.key)
		asserts.Nil(d.rbt().root.right)
		asserts.NotNil(d1.rbt().root.right)
		asserts.Equal(2, d1.rbt().root.right.key)
	})

	t.Run("Insert on Empty", func(t *testing.T) {
		d := Empty[int, int]()
		d1 := Insert(2, 2, d)

		asserts.Equal(&dict[int, int]{root: nil}, d.rbt())
		asserts.Equal(2, d1.rbt().root.key)
	})

	t.Run("Insert into existing entry", func(t *testing.T) {
		d := Singleton(10, 233)
		d1 := Insert(10, 233, d)

		SUT := d1.rbt()

		asserts.Equal(233, d.rbt().root.value)
		asserts.Equal(233, SUT.root.value)
		asserts.NotSame(d.rbt().root, SUT.root)
	})

	t.Run("Insert with color pushdown", func(t *testing.T) {
		d := Singleton(40, 1)
		d1 := Insert(50, 2, d)
		d2 := Insert(30, 3, d1)
		// Will cause color pushdown of parent node (40)
		d3 := Insert(35, 3, d2)

		asserts.Equal(d1.rbt().root.right, d2.rbt().root.right)
		asserts.NotEqual(d2.rbt().root.right, d3.rbt().root.right)
	})

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(&dict[int, struct{}]{root: nil}, Empty[int, struct{}]())
	})

	t.Run("Insert on Empty has properties", func(t *testing.T) {
		d := Empty[int, int]()
		d1 := Insert(1, 233, d)

		SUT := d1.rbt()

		asserts.Equal(&dict[int, int]{root: &node[int, int]{
			key:   1,
			value: 233,
			color: BLACK,
			left:  nil,
			right: nil}},
			SUT,
		)
	})

	t.Run("Insert on Singleton right side", func(t *testing.T) {
		d := Singleton[int, int](1, 1)
		d1 := Insert(2, 2, d)

		SUT := d1.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(RED, SUT.root.right.color)
	})

	t.Run("Insert into existing entry", func(t *testing.T) {
		d := Singleton(10, 233)
		d1 := Insert(10, 100, d)

		SUT := d1.rbt()

		asserts.Equal(&dict[int, int]{root: &node[int, int]{
			key:   10,
			value: 100,
			color: BLACK,
			left:  nil,
			right: nil},
		}, SUT)
	})

	t.Run("LL Single right rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		Insert(30, 3, d1)

		asserts.Nil(d1.rbt().root.right)
		asserts.Equal(40, d1.rbt().root.left.key)
	})

	t.Run("LR -> RR rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(45, 3, d1)

		// d1
		asserts.Equal(50, d1.rbt().root.key)
		asserts.Nil(d1.rbt().root.right)
		asserts.Equal(40, d1.rbt().root.left.key)
		asserts.Nil(d1.rbt().root.left.left)
		asserts.Nil(d1.rbt().root.left.right)

		// d2
		asserts.Equal(45, d2.rbt().root.key)
		asserts.Equal(50, d2.rbt().root.right.key)
		asserts.Nil(d2.rbt().root.right.right)
		asserts.Equal(40, d2.rbt().root.left.key)
		asserts.Nil(d2.rbt().root.left.left)
	})

	t.Run("RR Single left rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(60, 2, d)
		Insert(70, 3, d1)

		asserts.Nil(d1.rbt().root.left)
		asserts.Equal(60, d1.rbt().root.right.key)
	})

	t.Run("LR -> LL rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(60, 2, d)
		d2 := Insert(55, 3, d1)

		// d1
		asserts.Equal(50, d1.rbt().root.key)
		asserts.Nil(d1.rbt().root.left)
		asserts.Equal(60, d1.rbt().root.right.key)
		asserts.Nil(d1.rbt().root.right.right)
		asserts.Nil(d1.rbt().root.right.left)

		// d2
		asserts.Equal(55, d2.rbt().root.key)
		asserts.Equal(60, d2.rbt().root.right.key)
		asserts.Nil(d2.rbt().root.left.left)
		asserts.Equal(50, d2.rbt().root.left.key)
		asserts.Nil(d2.rbt().root.right.right)
	})

	t.Run("granparent color pushdown", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(45, 3, d1)
		d3 := Insert(30, 3, d2)

		// d2
		asserts.Equal(50, d2.rbt().root.right.key)
		asserts.Equal(RED, d2.rbt().root.right.color)
		asserts.Nil(d2.rbt().root.left.left)

		// d3
		asserts.Equal(50, d3.rbt().root.right.key)
		asserts.Equal(BLACK, d3.rbt().root.right.color)
		asserts.Equal(BLACK, d3.rbt().root.left.color)
		asserts.Equal(RED, d3.rbt().root.left.left.color)
	})

	t.Run("LL Single right rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(30, 3, d1)

		SUT := d2.rbt()

		asserts.Equal(40, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(50, SUT.root.right.key)
		asserts.Equal(RED, SUT.root.right.color)
		asserts.Equal(30, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.left.color)
	})

	t.Run("RR Single right rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(60, 2, d)
		d2 := Insert(70, 3, d1)

		SUT := d2.rbt()

		asserts.Equal(60, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(50, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(70, SUT.root.right.key)
		asserts.Equal(RED, SUT.root.right.color)
	})

	t.Run("LR double red, red uncle", func(t *testing.T) {
		d := Singleton(50, 1)
		// Left
		d1 := Insert(40, 2, d)
		d2 := Insert(60, 3, d1)
		d3 := Insert(45, 4, d2)

		SUT := d3.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(RED, SUT.root.left.right.color)
	})

	t.Run("LR double red, black uncle", func(t *testing.T) {
		d := Singleton(50, 1)
		// Left
		d1 := Insert(40, 2, d)
		d2 := Insert(45, 3, d1)

		SUT := d2.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(45, SUT.root.key)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(40, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.right.color)
		asserts.Equal(50, SUT.root.right.key)
	})

	t.Run("RL double red, red uncle", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(60, 2, d)
		d2 := Insert(40, 3, d1)
		d3 := Insert(55, 4, d2)

		SUT := d3.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(RED, SUT.root.right.left.color)
	})

	t.Run("RL double red, black uncle", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(60, 2, d)
		d2 := Insert(55, 4, d1)

		SUT := d2.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(55, SUT.root.key)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(50, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.right.color)
		asserts.Equal(60, SUT.root.right.key)
	})

	t.Run("test the following inserts 7,5,10,20,15", func(t *testing.T) {
		d := Singleton(7, 1)
		d1 := Insert(5, 2, d)
		d2 := Insert(10, 3, d1)
		d3 := Insert(20, 3, d2)
		d4 := Insert(15, 3, d3)

		SUT := d4.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(7, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(15, SUT.root.right.key)
		asserts.Equal(RED, SUT.root.right.right.color)
		asserts.Equal(20, SUT.root.right.right.key)
		asserts.Equal(RED, SUT.root.right.left.color)
		asserts.Equal(10, SUT.root.right.left.key)
	})

	t.Run("test the following inserts 10,15,5,0,2", func(t *testing.T) {
		d := Singleton(10, 1)
		d1 := Insert(15, 2, d)
		d2 := Insert(5, 3, d1)
		d3 := Insert(0, 3, d2)
		d4 := Insert(2, 3, d3)

		SUT := d4.rbt()

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(10, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(2, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.left.left.color)
		asserts.Equal(0, SUT.root.left.left.key)
		asserts.Equal(RED, SUT.root.left.right.color)
		asserts.Equal(5, SUT.root.left.right.key)
	})

	t.Run("Structure sharing right subtree", func(t *testing.T) {
		d := Singleton(40, 1)
		d1 := Insert(50, 2, d)
		d2 := Insert(30, 3, d1)
		d3 := Insert(35, 3, d1)

		asserts.Equal(d1.rbt().root.right, d2.rbt().root.right)
		asserts.Equal(d2.rbt().root.right, d3.rbt().root.right)
	})

}

func TestRemove(t *testing.T) {
	asserts := assert.New(t)

	t.Run("remove Empty Dict", func(t *testing.T) {
		d := Empty[int, int]()
		d1 := Remove(50, d)

		asserts.Equal(&d, &d1)
	})

	t.Run("remove Singleton key that doesn't exist", func(t *testing.T) {
		d := Singleton[int, int](1, 1)
		d1 := Remove(50, d)

		// Pointers match
		asserts.Equal(&d, &d1)
	})

	t.Run("remove Singleton", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Remove(50, d)

		asserts.NotNil(d.rbt().root)
		asserts.Equal(50, d.rbt().root.key)
		asserts.Nil(d1.rbt().root)
	})

	t.Run("remove childless RED leaf", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(60, 3, d1)
		d3 := Remove(40, d2)

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

	t.Run("BLACK leaf, LEFT | BLACK sibling, RED near nephew, BLACK distant nephew", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   40,
				value: 1,
				color: BLACK,
				left:  &node[int, int]{key: 30, value: 2, color: BLACK, left: nil, right: nil},
				right: &node[int, int]{
					key:   50,
					value: 3,
					color: BLACK,
					left:  &node[int, int]{key: 45, value: 4, color: RED, left: nil, right: nil},
					right: nil,
				},
			},
		}

		SUT := Remove(30, tree).rbt()

		asserts.Equal(40, tree.rbt().root.key)
		asserts.Equal(45, SUT.root.key)
		asserts.Equal(50, SUT.root.right.key)
		asserts.Equal(40, SUT.root.left.key)
		asserts.NotEqual(tree.rbt().root.right, SUT.root.right)
	})

	t.Run("BLACK leaf, RIGHT | BLACK sibling, RED near nephew, BLACK distant nephew", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   40,
				value: 1,
				color: BLACK,
				left: &node[int, int]{
					key:   30,
					value: 2,
					color: BLACK,
					left:  nil,
					right: &node[int, int]{key: 35, value: 4, color: RED, left: nil, right: nil},
				},
				right: &node[int, int]{
					key:   50,
					value: 3,
					color: BLACK,
					left:  nil,
					right: nil,
				},
			},
		}

		SUT := Remove(50, tree).rbt()

		asserts.Equal(40, tree.rbt().root.key)
		asserts.Equal(35, SUT.root.key)
		asserts.Equal(40, SUT.root.right.key)
		asserts.Equal(30, SUT.root.left.key)
		asserts.NotEqual(tree.rbt().root.left, SUT.root.left)
	})

	t.Run("BLACK leaf, RIGHT | RED sibling | BLACK near nephew | BLACK distant nephew", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   50,
				value: 1,
				color: BLACK,
				left: &node[int, int]{
					key:   40,
					value: 2,
					color: RED,
					left:  &node[int, int]{key: 35, value: 5, color: BLACK, left: nil, right: nil},
					right: &node[int, int]{
						key:   45,
						value: 6,
						color: BLACK,
						left:  nil,
						right: nil,
					},
				},
				right: &node[int, int]{
					key:   60,
					value: 3,
					color: BLACK,
					left:  nil,
					right: nil,
				},
			},
		}

		SUT := Remove(60, tree).rbt()

		// Removes node
		asserts.Nil(SUT.root.right.right)

		asserts.Equal(50, tree.rbt().root.key)
		asserts.Equal(40, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(50, SUT.root.right.key)
		asserts.Equal(45, SUT.root.right.left.key)
		asserts.Equal(RED, SUT.root.right.left.color)

		// Structure sharing
		asserts.True(tree.rbt().root.left.left == SUT.root.left)
	})

	t.Run("BLACK node, LEFT | RED child, LEFT | NIL child, RIGHT", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   50,
				value: 1,
				color: BLACK,
				left: &node[int, int]{
					key:   40,
					color: BLACK,
					value: 3,
					left:  nil,
					right: &node[int, int]{key: 45, color: RED, value: 6, left: nil, right: nil}},
				right: &node[int, int]{key: 60, color: BLACK, value: 2, left: nil, right: nil},
			},
		}

		SUT := Remove(40, tree).rbt()

		asserts.Nil(SUT.root.left.right)
		asserts.Equal(45, SUT.root.left.key)

		// Structure Sharing
		asserts.True(tree.rbt().root.right == SUT.root.right)
	})

	t.Run("BLACK node, RIGHT | RED child, LEFT | NIL child, RIGHT", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   50,
				value: 1,
				color: BLACK,
				right: &node[int, int]{
					key:   60,
					color: BLACK,
					value: 3,
					left:  &node[int, int]{key: 55, color: RED, value: 6, left: nil, right: nil},
					right: nil,
				},
				left: &node[int, int]{key: 40, color: BLACK, value: 2, left: nil, right: nil},
			},
		}

		SUT := Remove(60, tree).rbt()

		asserts.Nil(SUT.root.right.left)
		asserts.Equal(55, SUT.root.right.key)

		// Structure Sharing
		asserts.True(tree.rbt().root.left == SUT.root.left)
	})

	t.Run("Removes root node with 2 red children", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   50,
				color: BLACK,
				value: 1,
				left:  &node[int, int]{key: 40, color: RED, value: 2, left: nil, right: nil},
				right: &node[int, int]{key: 60, color: RED, value: 3, left: nil, right: nil},
			},
		}

		SUT := Remove(50, tree).rbt()

		asserts.Equal(60, SUT.root.key)
		asserts.Nil(SUT.root.right)
		asserts.Equal(40, SUT.root.left.key)

		// Structure Sharing
		asserts.Equal(tree.rbt().root.left, SUT.root.left)
		asserts.NotEqual(tree.rbt().root, SUT.root)
	})

	t.Run("Removes red right leaf node with no children", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(60, 3, d1)
		d3 := Remove(60, d2)

		SUT := d3.rbt()

		asserts.Nil(SUT.root.right)

		// Structure Sharing
		asserts.True(d1.rbt().root.left == SUT.root.left)
		asserts.True(d1.rbt().root != SUT.root)
	})

	t.Run("Removes a red left node with no children", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(60, 3, d1)
		d3 := Remove(40, d2)

		SUT := d3.rbt()

		asserts.Nil(SUT.root.left)
	})

	t.Run("Solve rbt | remove 50,20,100,90,40,60,70,10,30,80", func(t *testing.T) {
		// https://www.youtube.com/watch?v=PgO_Xj7DC1A&t=16s
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   40,
				value: 1,
				color: BLACK,
				left: &node[int, int]{
					key:   20,
					value: 2,
					color: BLACK,
					left:  &node[int, int]{key: 10, value: 3, color: BLACK, left: nil, right: nil},
					right: &node[int, int]{key: 30, value: 4, color: BLACK, left: nil, right: nil},
				},
				right: &node[int, int]{
					key:   60,
					value: 5,
					color: BLACK,
					left:  &node[int, int]{key: 50, value: 6, color: BLACK, left: nil, right: nil},
					right: &node[int, int]{
						key:   80,
						value: 7,
						color: RED,
						left:  &node[int, int]{key: 70, value: 8, color: BLACK, left: nil, right: nil},
						right: &node[int, int]{
							key:   90,
							value: 9,
							color: BLACK,
							left:  nil,
							right: &node[int, int]{key: 100, value: 10, color: RED, left: nil, right: nil},
						},
					},
				},
			},
		}

		// REMOVE 50
		tree1 := Remove(50, tree)
		SUT1 := tree1.rbt()

		asserts.Equal(40, SUT1.rbt().root.key)
		asserts.Equal(BLACK, SUT1.rbt().root.color)

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
		tree2 := Remove(20, tree1)
		SUT2 := tree2.rbt()

		asserts.Equal(40, SUT2.root.key)
		asserts.Equal(BLACK, SUT2.root.color)

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
		tree3 := Remove(100, tree2)
		SUT3 := tree3.rbt()

		// Different root
		asserts.True(SUT2.root != SUT3.root)

		// Struture sharing
		// 50
		asserts.True(SUT2.root.left == SUT3.root.left)
		// 10
		asserts.True(SUT2.root.left.left == SUT3.root.left.left)

		// REMOVE 90
		tree4 := Remove(90, tree3)
		SUT4 := tree4.rbt()

		asserts.Equal(40, SUT4.root.key)
		asserts.Equal(70, SUT4.root.right.key)
		asserts.Equal(RED, SUT4.root.right.color)
		asserts.Equal(80, SUT4.root.right.right.key)
		asserts.Equal(BLACK, SUT4.root.right.right.color)
		asserts.Equal(60, SUT4.root.right.left.key)
		asserts.Equal(BLACK, SUT4.root.right.left.color)
		asserts.Nil(SUT4.root.right.left.left)
		asserts.Nil(SUT4.root.right.left.right)

		// REMOVE 40
		tree5 := Remove(40, tree4)
		SUT5 := tree5.rbt()

		asserts.Equal(60, SUT5.root.key)
		asserts.Equal(BLACK, SUT5.root.color)
		asserts.Equal(70, SUT5.root.right.key)
		asserts.Equal(BLACK, SUT5.root.right.color)
		asserts.Equal(80, SUT5.root.right.right.key)
		asserts.Equal(RED, SUT5.root.right.right.color)

		// Structure sharing
		asserts.True(tree5.rbt().root.left == SUT5.root.left)
		asserts.True(tree5.rbt().root.left.left == SUT5.root.left.left)

		// REMOVE 60
		tree6 := Remove(60, tree5)
		SUT6 := tree6.rbt()

		asserts.Equal(70, SUT6.root.key)
		asserts.Equal(BLACK, SUT6.root.color)
		asserts.Equal(80, SUT6.root.right.key)
		asserts.Equal(BLACK, SUT6.root.right.color)

		// REMOVE 70
		tree7 := Remove(70, tree6)
		SUT7 := tree7.rbt()

		asserts.Equal(30, SUT7.root.key)
		asserts.Equal(BLACK, SUT7.root.color)
		asserts.Equal(80, SUT7.root.right.key)
		asserts.Equal(BLACK, SUT7.root.right.color)

		// Structure sharing
		asserts.True(tree7.rbt().root.left == SUT7.root.left)

		// REMOVE 10
		tree8 := Remove(10, tree7)
		SUT8 := tree8.rbt()

		asserts.Equal(30, SUT8.root.key)
		asserts.Equal(BLACK, SUT8.root.color)
		asserts.Equal(80, SUT8.root.right.key)
		asserts.Equal(RED, SUT8.root.right.color)
		asserts.Nil(SUT8.root.left)

		// REMOVE 30
		tree9 := Remove(30, tree8)
		SUT9 := tree9.rbt()

		asserts.Equal(80, SUT9.root.key)
		asserts.Equal(BLACK, SUT9.root.color)
		asserts.Nil(SUT9.root.right)

		// REMOVE 80
		tree10 := Remove(80, tree9)
		SUT10 := tree10.rbt()

		asserts.Nil(SUT10.root)
	})
}
