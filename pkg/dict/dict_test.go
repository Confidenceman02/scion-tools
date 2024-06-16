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

	t.Run("Immutable persistent Insert of Singleton", func(t *testing.T) {
		d := Singleton(1, 2)
		d1 := Insert(2, 2, d)

		asserts.Equal(1, d.rbt().root.key)
		asserts.Equal(1, d1.rbt().root.key)
		asserts.Nil(d.rbt().root.right)
		asserts.NotNil(d1.rbt().root.right)
		asserts.Equal(2, d1.rbt().root.right.key)
	})

	t.Run("Immutable persistent Insert of Empty", func(t *testing.T) {
		d := Empty[int, int]()
		d1 := Insert(2, 2, d)

		asserts.Equal(&dict[int, int]{root: nil}, d.rbt())
		asserts.Equal(2, d1.rbt().root.key)
	})

	t.Run("Immutable persistent Insert into existing entry", func(t *testing.T) {
		d := Singleton(10, 233)
		d1 := Insert(10, 233, d)

		SUT := d1.rbt()

		asserts.Equal(233, d.rbt().root.value)
		asserts.Equal(233, SUT.root.value)
		asserts.NotSame(d.rbt().root, SUT.root)
	})

	t.Run("Immutable persistent Insert with color pushdown", func(t *testing.T) {
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

	t.Run("Immutable persistent LL Single right rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		Insert(30, 3, d1)

		asserts.Nil(d1.rbt().root.right)
		asserts.Equal(40, d1.rbt().root.left.key)
	})

	t.Run("Immutable persistent LR - RR rotation", func(t *testing.T) {
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

	t.Run("Immutable persistent RR Single left rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(60, 2, d)
		Insert(70, 3, d1)

		asserts.Nil(d1.rbt().root.left)
		asserts.Equal(60, d1.rbt().root.right.key)
	})

	t.Run("Immutable persistent LR - LL rotation", func(t *testing.T) {
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

	t.Run("Immutable persistent after granparent color pushdown", func(t *testing.T) {
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
}

func TestRemove(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Immutable persistent removal of Empty Dict", func(t *testing.T) {
		d := Empty[int, int]()
		d1 := Remove(50, d)

		asserts.Equal(&d, &d1)
	})

	t.Run("Immutable persistent removal of Singleton key that doesn't exist", func(t *testing.T) {
		d := Singleton[int, int](1, 1)
		d1 := Remove(50, d)

		// Pointers match
		asserts.Equal(&d, &d1)
	})

	t.Run("Immutable persistent removal of Singleton", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Remove(50, d)

		asserts.NotNil(d.rbt().root)
		asserts.Equal(50, d.rbt().root.key)
		asserts.Nil(d1.rbt().root)
	})

	t.Run("Immutable persistent removal of childless RED leaf", func(t *testing.T) {
		d := Singleton(50, 1)
		d1 := Insert(40, 2, d)
		d2 := Insert(60, 3, d1)
		d3 := Remove(40, d2)

		asserts.NotNil(d1.rbt().root.left)
		asserts.NotNil(d2.rbt().root.left)
		asserts.Nil(d3.rbt().root.left)
		// Shared structure d2 - d3
		asserts.Equal(d2.rbt().root.right, d3.rbt().root.right)
		// Different root d2, d3
		asserts.NotEqual(d2.rbt().root, d3.rbt().root)
	})

	t.Run("Immutable persistent removal of | BLACK leaf, LEFT | BLACK sibling, RED near nephew, BLACK distant nephew", func(t *testing.T) {
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

	t.Run("Immutable persistent removal of | BLACK leaf, RIGHT | BLACK sibling, RED near nephew, BLACK distant nephew", func(t *testing.T) {
		var tree Dict[int, int]
		tree = &dict[int, int]{
			root: &node[int, int]{
				key:   40,
				value: 1,
				color: BLACK,
				left:  &node[int, int]{key: 30, value: 2, color: BLACK, left: nil, right: &node[int, int]{key: 35, value: 4, color: RED, left: nil, right: nil}},
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

	t.Run("Immutable persistent removal of | BLACK leaf, RIGHT | RED sibling", func(t *testing.T) {
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

		asserts.Equal(50, tree.rbt().root.key)
		asserts.Equal(40, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(50, SUT.root.right.key)
		asserts.Equal(45, SUT.root.right.left.key)
		asserts.Equal(RED, SUT.root.right.left.color)
		asserts.Equal(tree.rbt().root.left.left, SUT.root.left)
	})

	// t.Run("Immutable persistent removal of childless BLACK leaf", func(t *testing.T) {
	// 	d := Singleton(40, 1)
	// 	d1 := Insert(50, 2, d)
	// 	d2 := Insert(30, 3, d1)
	// 	d3 := Insert(35, 3, d1)
	//
	// 	asserts.NotNil(d1.rbt().root.left)
	// 	asserts.NotNil(d2.rbt().root.left)
	// 	asserts.Nil(d3.rbt().root.left)
	// 	// Shared structure d2 - d3
	// 	asserts.Equal(d2.rbt().root.right, d3.rbt().root.right)
	// 	// Different root d2, d3
	// 	asserts.NotEqual(d2.rbt().root, d3.rbt().root)
	// })

	// t.Run("Removes root node with 2 red children", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(60, 2)
	// 	d2 := d1.Insert(40, 3)
	// 	d3 := d2.Remove(50)
	//
	// 	SUT := d3.rbt()
	//
	// 	asserts.Equal(60, SUT.root.key)
	// 	asserts.Nil(SUT.root.right)
	// 	asserts.Equal(40, SUT.root.left.key)
	// })
	//
	// t.Run("Removes red right leaf node with no children", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(40, 2)
	// 	d2 := d1.Insert(60, 3)
	// 	d3 := d2.Remove(60)
	//
	// 	SUT := d3.rbt()
	//
	// 	asserts.Nil(SUT.root.right)
	// })
	//
	// t.Run("Removes a red left node with no children", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(40, 2)
	// 	d2 := d1.Insert(60, 3)
	// 	d3 := d2.Remove(40)
	//
	// 	SUT := d3.rbt()
	//
	// 	asserts.Nil(SUT.root.left)
	// })
	//
	// t.Run("Removes a black leaf node with 1 child | Left", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(40, 2)
	// 	d2 := d1.Insert(60, 3)
	// 	d3 := d2.Insert(45, 4)
	//
	// 	SUT := d3.rbt()
	//
	// 	SUT.Remove(40)
	//
	// 	asserts.Equal(50, SUT.root.key)
	// 	asserts.Equal(BLACK, SUT.root.color)
	// 	asserts.Equal(45, SUT.root.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.color)
	// 	asserts.Nil(SUT.root.left.right)
	// })
	//
	// t.Run("Removes a black leaf node with 1 child | Right", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(40, 2)
	// 	d2 := d1.Insert(60, 3)
	// 	d3 := d2.Insert(55, 4)
	//
	// 	SUT := d3.rbt()
	//
	// 	SUT.Remove(60)
	//
	// 	asserts.Equal(50, SUT.root.key)
	// 	asserts.Equal(BLACK, SUT.root.color)
	// 	asserts.Equal(55, SUT.root.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Nil(SUT.root.right.left)
	// })
	//
	// t.Run("Removes a black leaf node with no children | p = RED | s = BLACK with no children", func(t *testing.T) {
	// 	d := Singleton(10, 1)
	// 	d1 := d.Insert(5, 2)
	// 	d2 := d1.Insert(20, 3)
	// 	d3 := d2.Insert(15, 4)
	// 	d4 := d3.Insert(30, 5)
	//
	// 	d5 := d4.rbt()
	//
	// 	// Mutate tree to for testing
	// 	d5.root.right.color = RED
	// 	d5.root.right.right.color = BLACK
	// 	d5.root.right.left.color = BLACK
	// 	d5.root.left.color = BLACK
	//
	// 	SUT := d5.Remove(15).rbt()
	//
	// 	asserts.Nil(SUT.root.right.left)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Equal(RED, SUT.root.right.right.color)
	// })
	//
	// t.Run("Removes a black leaf node with no children | p = BLACK | s = BLACK with no children", func(t *testing.T) {
	// 	d := Singleton(10, 1)
	// 	d1 := d.Insert(5, 2)
	// 	d2 := d1.Insert(20, 3)
	// 	d3 := d2.Insert(1, 2)
	// 	d4 := d3.Insert(7, 2)
	// 	d5 := d4.Insert(15, 4)
	// 	d6 := d5.Insert(30, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Manually balance for testing scenario
	// 	// RIGHT
	// 	d7.root.right.color = BLACK
	// 	d7.root.right.right.color = BLACK
	// 	d7.root.right.left.color = BLACK
	// 	// LEFT
	// 	d7.root.left.color = BLACK
	// 	d7.root.left.left.color = BLACK
	// 	d7.root.left.right.color = BLACK
	//
	// 	SUT := d7.Remove(15).rbt()
	//
	// 	asserts.Nil(SUT.root.right.left)
	// 	asserts.Equal(BLACK, SUT.root.color)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Equal(RED, SUT.root.left.color)
	// 	asserts.Equal(RED, SUT.root.right.right.color)
	// })
	//
	// t.Run("Removes a black leaf node with no children | p = BLACK | s = RED | right branch", func(t *testing.T) {
	// 	d := Singleton(10, 1)
	// 	d1 := d.Insert(5, 2)
	// 	d2 := d1.Insert(20, 3)
	// 	d3 := d2.Insert(1, 2)
	// 	d4 := d3.Insert(7, 2)
	// 	d5 := d4.Insert(15, 4)
	// 	d6 := d5.Insert(30, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Mutate tree
	// 	// RIGHT
	// 	d7.root.right.color = BLACK
	// 	d7.root.right.right.color = RED
	// 	d7.root.right.left.color = BLACK
	// 	// LEFT
	// 	d7.root.left.color = BLACK
	// 	d7.root.left.left.color = BLACK
	// 	d7.root.left.right.color = BLACK
	//
	// 	// Balance
	// 	d7.root.right.right.right = &node[int, int]{parent: d7.root.right.right, key: 40, value: 6, color: BLACK}
	// 	d7.root.right.right.left = &node[int, int]{parent: d7.root.right.right, key: 25, value: 7, color: BLACK}
	//
	// 	SUT := d7.Remove(15).rbt()
	//
	// 	asserts.Equal(30, SUT.root.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Equal(40, SUT.root.right.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.right.color)
	// 	asserts.Equal(20, SUT.root.right.left.key)
	// 	asserts.Equal(BLACK, SUT.root.right.left.color)
	// 	asserts.Equal(25, SUT.root.right.left.right.key)
	// 	asserts.Equal(RED, SUT.root.right.left.right.color)
	// 	asserts.Nil(SUT.root.right.left.left)
	// })
	//
	// t.Run("Removes a black leaf node with no children | p = BLACK | s = RED | left branch", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(40, 2)
	// 	d2 := d1.Insert(60, 3)
	// 	d3 := d2.Insert(70, 2)
	// 	d4 := d3.Insert(55, 2)
	// 	d5 := d4.Insert(45, 4)
	// 	d6 := d5.Insert(35, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Mutate tree
	// 	// LEFT
	// 	d7.root.left.color = BLACK
	// 	d7.root.left.left.color = RED
	// 	d7.root.left.right.color = BLACK
	// 	// RIGHT
	// 	d7.root.right.color = BLACK
	// 	d7.root.right.right.color = BLACK
	// 	d7.root.right.left.color = BLACK
	//
	// 	// Balance
	// 	d7.root.left.left.left = &node[int, int]{parent: d7.root.left.left, key: 20, value: 6, color: BLACK}
	// 	d7.root.left.left.right = &node[int, int]{parent: d7.root.left.left, key: 37, value: 7, color: BLACK}
	//
	// 	SUT := d7.Remove(45).rbt()
	//
	// 	asserts.Equal(35, SUT.root.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.color)
	// 	asserts.Equal(20, SUT.root.left.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.left.color)
	// 	asserts.Equal(40, SUT.root.left.right.key)
	// 	asserts.Equal(BLACK, SUT.root.left.right.color)
	// 	asserts.Equal(37, SUT.root.left.right.left.key)
	// 	asserts.Equal(RED, SUT.root.left.right.left.color)
	// 	asserts.Nil(SUT.root.left.right.right)
	// })
	//
	// t.Run("DB | s = BLACK with red and black child | Left subtree", func(t *testing.T) {
	// 	// From example https://www.youtube.com/watch?v=4KDovab_OS8&list=PLmp4WtRF6yg0_07IUb2eOsS0k1jIa2IgP&index=5&t=1819s
	// 	d := Singleton(10, 1)
	// 	d1 := d.Insert(5, 2)
	// 	d2 := d1.Insert(30, 3)
	// 	d3 := d2.Insert(25, 2)
	// 	d4 := d3.Insert(40, 2)
	// 	d5 := d4.Insert(7, 4)
	// 	d6 := d5.Insert(1, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Mutate tree for example
	// 	// LEFT
	// 	d7.root.left.left.color = BLACK
	// 	d7.root.left.right.color = BLACK
	// 	// RIGHT
	// 	d7.root.right.right.color = BLACK
	// 	d7.root.right.left.color = RED
	//
	// 	// Manually Balance
	// 	d7.root.right.left.left = &node[int, int]{parent: d7.root.right.left, key: 20, value: 6, color: BLACK}
	// 	d7.root.right.left.right = &node[int, int]{parent: d7.root.right.left, key: 28, value: 7, color: BLACK}
	//
	// 	SUT := d7.Remove(1).rbt()
	//
	// 	asserts.Equal(25, SUT.root.key)
	// 	asserts.Equal(10, SUT.root.left.key)
	// 	asserts.Equal(30, SUT.root.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Equal(5, SUT.root.left.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.left.color)
	// 	asserts.Equal(7, SUT.root.left.left.right.key)
	// 	asserts.Equal(RED, SUT.root.left.left.right.color)
	// 	asserts.Equal(20, SUT.root.left.right.key)
	// })
	//
	// t.Run("DB | s = BLACK with red and black child | Right subtree", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(60, 2)
	// 	d2 := d1.Insert(40, 3)
	// 	d3 := d2.Insert(45, 2)
	// 	d4 := d3.Insert(30, 2)
	// 	d5 := d4.Insert(55, 4)
	// 	d6 := d5.Insert(70, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Mutate tree for testing
	// 	// LEFT
	// 	d7.root.left.left.color = BLACK
	// 	// RIGHT
	// 	d7.root.right.right.color = BLACK
	// 	d7.root.right.left.color = BLACK
	//
	// 	// Manually Balance
	// 	d7.root.left.right.right = &node[int, int]{parent: d7.root.left.right, key: 47, value: 6, color: BLACK}
	// 	d7.root.left.right.left = &node[int, int]{parent: d7.root.left.right, key: 41, value: 7, color: BLACK}
	//
	// 	SUT := d7.Remove(70).rbt()
	//
	// 	asserts.Equal(45, SUT.root.key)
	// 	asserts.Equal(BLACK, SUT.root.color)
	// 	asserts.Equal(40, SUT.root.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.color)
	// 	asserts.Equal(50, SUT.root.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Equal(30, SUT.root.left.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.left.color)
	// 	asserts.Equal(60, SUT.root.right.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.right.color)
	// 	asserts.Equal(47, SUT.root.right.left.key)
	// 	asserts.Equal(BLACK, SUT.root.right.left.color)
	// 	asserts.Equal(55, SUT.root.right.right.left.key)
	// 	asserts.Equal(RED, SUT.root.right.right.left.color)
	// 	asserts.Nil(SUT.root.right.right.right)
	// })
	//
	// t.Run("DB | s = BLACK with red and black child | Right subtree", func(t *testing.T) {
	// 	d := Singleton(50, 1)
	// 	d1 := d.Insert(60, 2)
	// 	d2 := d1.Insert(40, 3)
	// 	d3 := d2.Insert(45, 2)
	// 	d4 := d3.Insert(30, 2)
	// 	d5 := d4.Insert(55, 4)
	// 	d6 := d5.Insert(70, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Mutate tree for testing
	// 	// LEFT
	// 	d7.root.left.left.color = BLACK
	// 	// RIGHT
	// 	d7.root.right.right.color = BLACK
	// 	d7.root.right.left.color = BLACK
	//
	// 	// Manually Balance
	// 	d7.root.left.right.right = &node[int, int]{parent: d7.root.left.right, key: 47, value: 6, color: BLACK}
	// 	d7.root.left.right.left = &node[int, int]{parent: d7.root.left.right, key: 41, value: 7, color: BLACK}
	//
	// 	SUT := d7.Remove(70).rbt()
	//
	// 	asserts.Equal(45, SUT.root.key)
	// 	asserts.Equal(BLACK, SUT.root.color)
	// 	asserts.Equal(40, SUT.root.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.color)
	// 	asserts.Equal(50, SUT.root.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.color)
	// 	asserts.Equal(30, SUT.root.left.left.key)
	// 	asserts.Equal(BLACK, SUT.root.left.left.color)
	// 	asserts.Equal(60, SUT.root.right.right.key)
	// 	asserts.Equal(BLACK, SUT.root.right.right.color)
	// 	asserts.Equal(47, SUT.root.right.left.key)
	// 	asserts.Equal(BLACK, SUT.root.right.left.color)
	// 	asserts.Equal(55, SUT.root.right.right.left.key)
	// 	asserts.Equal(RED, SUT.root.right.right.left.color)
	// 	asserts.Nil(SUT.root.right.right.right)
	// })
	//
	// t.Run("Solve rbt | remove 50,20,100,90,40,60,70,10,30,80", func(t *testing.T) {
	// 	// https://www.youtube.com/watch?v=PgO_Xj7DC1A&t=16s
	// 	d := Singleton(40, 1)
	// 	d1 := d.Insert(20, 2)
	// 	d2 := d1.Insert(60, 3)
	// 	d3 := d2.Insert(10, 2)
	// 	d4 := d3.Insert(30, 2)
	// 	d5 := d4.Insert(50, 4)
	// 	d6 := d5.Insert(80, 5)
	//
	// 	d7 := d6.rbt()
	//
	// 	// Mutate tree for testing
	// 	// LEFT
	// 	d7.root.left.left.color = BLACK
	// 	d7.root.left.right.color = BLACK
	// 	// RIGHT
	// 	d7.root.right.left.color = BLACK
	// 	// SUT.root.right.right.color = RED
	//
	// 	// Manually Balance
	// 	d7.root.right.right.left = &node[int, int]{parent: d7.root.right.right, key: 70, value: 6, color: BLACK}
	// 	d7.root.right.right.right = &node[int, int]{parent: d7.root.right.right, key: 90, value: 7, color: BLACK}
	// 	d7.root.right.right.right.right = &node[int, int]{parent: d7.root.right.right.right, key: 100, value: 7, color: RED}
	//
	// 	// REMOVE 50
	// 	SUT1 := d7.Remove(50).rbt()
	//
	// 	asserts.Equal(40, SUT1.root.key)
	// 	asserts.Equal(BLACK, SUT1.root.color)
	// 	// right
	// 	asserts.Equal(80, SUT1.root.right.key)
	// 	asserts.Equal(BLACK, SUT1.root.right.color)
	// 	// left
	// 	asserts.Equal(20, SUT1.root.left.key)
	// 	asserts.Equal(BLACK, SUT1.root.left.color)
	// 	// right right
	// 	asserts.Equal(90, SUT1.root.right.right.key)
	// 	asserts.Equal(BLACK, SUT1.root.right.right.color)
	// 	// left left
	// 	asserts.Equal(10, SUT1.root.left.left.key)
	// 	asserts.Equal(BLACK, SUT1.root.left.left.color)
	// 	// right right right
	// 	asserts.Equal(100, SUT1.root.right.right.right.key)
	// 	asserts.Equal(RED, SUT1.root.right.right.right.color)
	// 	// right left
	// 	asserts.Equal(60, SUT1.root.right.left.key)
	// 	asserts.Equal(BLACK, SUT1.root.right.left.color)
	// 	// left right
	// 	asserts.Equal(30, SUT1.root.left.right.key)
	// 	asserts.Equal(BLACK, SUT1.root.left.right.color)
	// 	// right left right
	// 	asserts.Equal(70, SUT1.root.right.left.right.key)
	// 	asserts.Equal(RED, SUT1.root.right.left.right.color)
	//
	// 	asserts.Nil(SUT1.root.right.left.right.left)
	// 	asserts.Nil(SUT1.root.right.left.right.right)
	//
	// 	// REMOVE 20
	// 	SUT2 := SUT1.Remove(20).rbt()
	//
	// 	asserts.Equal(40, SUT2.root.key)
	// 	asserts.Equal(BLACK, SUT2.root.color)
	// 	// left
	// 	asserts.Equal(30, SUT2.root.left.key)
	// 	asserts.Equal(BLACK, SUT2.root.left.color)
	// 	// left left
	// 	asserts.Equal(10, SUT2.root.left.left.key)
	// 	asserts.Equal(RED, SUT2.root.left.left.color)
	// 	// left right
	// 	asserts.Nil(SUT2.root.left.right)
	// 	// right
	// 	asserts.Equal(80, SUT2.root.right.key)
	// 	asserts.Equal(RED, SUT2.root.right.color)
	//
	// 	// REMOVE 100
	// 	SUT3 := SUT2.Remove(100).rbt()
	//
	// 	// right right right
	// 	asserts.Nil(SUT3.root.right.right.right)
	//
	// 	// REMOVE 90
	// 	SUT3.Remove(90)
	//
	// 	asserts.Equal(40, SUT3.root.key)
	// 	// right
	// 	asserts.Equal(70, SUT3.root.right.key)
	// 	asserts.Equal(RED, SUT3.root.right.color)
	// 	// right right
	// 	asserts.Equal(80, SUT3.root.right.right.key)
	// 	asserts.Equal(BLACK, SUT3.root.right.right.color)
	// 	// right right right
	// 	asserts.Nil(SUT3.root.right.right.right)
	// 	// right left
	// 	asserts.Equal(60, SUT3.root.right.left.key)
	// 	asserts.Equal(BLACK, SUT3.root.right.left.color)
	// 	// right left left
	// 	asserts.Nil(SUT3.root.right.left.left)
	// 	// right left right
	// 	asserts.Nil(SUT3.root.right.left.right)
	//
	// 	// REMOVE 40
	// 	SUT4 := SUT3.Remove(40).rbt()
	//
	// 	asserts.Equal(60, SUT4.root.key)
	// 	asserts.Equal(BLACK, SUT4.root.color)
	// 	// right
	// 	asserts.Equal(70, SUT4.root.right.key)
	// 	asserts.Equal(BLACK, SUT4.root.right.color)
	// 	// right right
	// 	asserts.Equal(80, SUT4.root.right.right.key)
	// 	asserts.Equal(RED, SUT4.root.right.right.color)
	//
	// 	// REMOVE 60
	// 	SUT5 := SUT4.Remove(60).rbt()
	//
	// 	asserts.Equal(70, SUT5.root.key)
	// 	asserts.Equal(BLACK, SUT5.root.color)
	// 	// // right
	// 	asserts.Equal(80, SUT5.root.right.key)
	// 	asserts.Equal(BLACK, SUT5.root.right.color)
	//
	// 	// REMOVE 70
	// 	SUT6 := SUT5.Remove(70).rbt()
	//
	// 	asserts.Equal(30, SUT6.root.key)
	// 	asserts.Equal(BLACK, SUT6.root.color)
	// 	// right
	// 	asserts.Equal(80, SUT6.root.right.key)
	// 	asserts.Equal(BLACK, SUT6.root.right.color)
	// 	// left
	// 	asserts.Equal(10, SUT6.root.left.key)
	// 	asserts.Equal(BLACK, SUT6.root.left.color)
	//
	// 	// REMOVE 10
	// 	SUT7 := SUT6.Remove(10).rbt()
	//
	// 	asserts.Equal(30, SUT7.root.key)
	// 	asserts.Equal(BLACK, SUT7.root.color)
	// 	// right
	// 	asserts.Equal(80, SUT7.root.right.key)
	// 	asserts.Equal(RED, SUT7.root.right.color)
	// 	// left
	// 	asserts.Nil(SUT7.root.left)
	//
	// 	// REMOVE 30
	// 	SUT8 := SUT7.Remove(30).rbt()
	//
	// 	asserts.Equal(80, SUT8.root.key)
	// 	asserts.Equal(BLACK, SUT8.root.color)
	// 	// right
	// 	asserts.Nil(SUT8.root.right)
	// 	// left
	// 	asserts.Nil(SUT8.root.left)
	//
	// 	// REMOVE 80
	// 	SUT9 := SUT8.Remove(80).rbt()
	//
	// 	asserts.Nil(SUT9.root)
	// })
}
