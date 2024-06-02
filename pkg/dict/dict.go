package dict

import (
	"cmp"
	"errors"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
)

const (
	EQ = 0
	LT = -1
	GT = +1
)

const (
	LEFT = iota
	RIGHT
)

const (
	RED int = iota
	BLACK
)

/*
Elm Dicts under the hood are Red-Black trees.

Insertions are standard BST insertions.

The rules are as follows for Red-Black trees:

 1. Every node is colored red or black.
 2. Every leaf is a NIL node, and is colored black.
 3. If a node is red, then both its children are black. (No consecutive red nodes)
 4. Every simple path from a node to a descendant leaf contains the same number of black nodes.
*/
type Dict[K cmp.Ordered, V any] interface {
	rbt() *dict[K, V]
}

/*
Retrieve the internal Red-Black tree
*/
func (d dict[K, V]) rbt() *dict[K, V] {
	return &d
}

type dict[K cmp.Ordered, V any] struct {
	root *node[K, V]
}

type node[K cmp.Ordered, V any] struct {
	key   K
	value V
	color int
	left  *node[K, V]
	right *node[K, V]
}

type stack[K cmp.Ordered, V any] struct {
	pp *stack[K, V]
	p  *node[K, V]
}

// Builders
func Empty[K cmp.Ordered, V any]() Dict[K, V] {
	return dict[K, V]{root: nil}
}

func Singleton[K cmp.Ordered, V any](key K, value V) Dict[K, V] {
	// Root nodes are always black
	return dict[K, V]{
		root: &node[K, V]{
			key:   key,
			value: value,
			color: BLACK,
			left:  nil,
			right: nil},
	}
}

// Methods
func Get[K cmp.Ordered, V any](targetKey K, d Dict[K, V]) maybe.Maybe[V] {
	root := d.rbt().root
	if root == nil {
		return maybe.Nothing{}
	} else {
		return getHelp(targetKey, root)
	}
}

func getHelp[K cmp.Ordered, V any](targetKey K, n *node[K, V]) maybe.Maybe[V] {
	if n != nil {
		switch cmp.Compare(targetKey, n.key) {
		case LT:
			return getHelp(targetKey, n.left)
		case EQ:
			return maybe.Just[V]{Value: n.value}
		case GT:
			return getHelp(targetKey, n.right)
		}
	}
	return maybe.Nothing{}
}

func (d *dict[K, V]) getNode(targetKey K) *node[K, V] {
	if d.root != nil {
		return getNodeHelp(targetKey, d.root)
	}
	return d.root
}

func getNodeHelp[K cmp.Ordered, V any](targetKey K, n *node[K, V]) *node[K, V] {
	if n != nil {
		switch cmp.Compare(targetKey, n.key) {
		case LT:
			if n.left == nil {
				return nil
			}
			return getNodeHelp(targetKey, n.left)
		case EQ:
			return n
		case GT:
			if n.right == nil {
				return nil
			}
			return getNodeHelp(targetKey, n.right)
		}
	}
	return nil
}

/*
Insert with BST insertion

- New nodes inserted are always red

Case 1 - Node is root
    1.1 Color node Black and exit

Case 2 - Black parent
    2.1 Exit

Case 3 - Parent is red and uncle is red
    3.1 Push down blackness from grandparent
    3.2 Find new condition for grandparent

Case 4 - Parent is red and uncle is Black
    LL
        ll.1 Rotate grandparent right
        ll.2 Swap colors of grandparent and parent
    LR
        LR.1 Left rotation of parent
        LR.2 Apply LL
    RR
        RR.1 Rotate grandparent left
        RR.2 Swap colors of grandparent and parent
    RL
        RL.1 Right rotation of parent
        RL.2 Apply RR
*/

func Insert[K cmp.Ordered, V any](key K, v V, d Dict[K, V]) Dict[K, V] {
	root := d.rbt().root

	if root == nil {
		return &dict[K, V]{root: &node[K, V]{key: key, value: v, color: BLACK, left: nil, right: nil}}
	}

	valRoot := *root
	inserted, stk := insertHelp(key, v, &valRoot, &stack[K, V]{pp: nil, p: nil})
	n, stk := balance(inserted, stk)
	newRoot := getStackRoot(n, stk)
	return &dict[K, V]{root: newRoot}
}

func insertHelp[K cmp.Ordered, V any](key K, value V, n *node[K, V], stk *stack[K, V]) (*node[K, V], *stack[K, V]) {
	nKey := n.key
	switch cmp.Compare(key, nKey) {
	case LT:
		if n.left == nil {
			n.left = &node[K, V]{key: key, value: value, color: RED, left: nil, right: nil}
			newStk := &stack[K, V]{pp: stk}
			newStk.p = n
			return n.left, newStk
		} else {
			newStk := &stack[K, V]{pp: stk}
			valL := *n.left
			n.left = &valL
			newStk.p = n
			return insertHelp(key, value, n.left, newStk)
		}
	case EQ:
		n.value = value
		return n, stk
	case GT:
		if n.right == nil {
			newStk := &stack[K, V]{pp: stk}
			newStk.p = n
			n.right = &node[K, V]{key: key, value: value, color: RED, left: nil, right: nil}
			return n.right, newStk
		} else {
			newStk := &stack[K, V]{pp: stk}
			valR := *n.right
			n.right = &valR
			newStk.p = n
			return insertHelp(key, value, n.right, newStk)
		}
	}
	panic("unreachable")
}

/*
Removal is a bit more of a process to that of insertion

For any node, the black height across all paths is equal.

- Only leaf nodes can be removed

Case 1 - Node is a red leaf
    1.1
        Delete node and exit
Case 2 - Double Black (DB) is root
    2.2
        Remove DB
Case 3 - DB sibling is black and both nephews are black
    3.1
        Remove DB node
    3.2
        Make sibling red
    3.3
        Add black to parent. If parent was red, make black
        otherwise make it a DB and find appropriate CASE
Case 4 - DB sibling is red
    4.1 Swap colors of DB parent & sibling
    4.2 Rotate parent in DB's direction
    4.3 Find next case for DB
Case 5 - DB sibling is black, far nephew is black and near nephew is red
    5.1 Swap colors of the DB sibling and near nephew
    5.2 Rotate sibling of DB node in opposite direction of DB node
    5.3 Apply case 6
Case 6 - DB sibling is black and far nephew is red
    6.1 Swap the colors of the DB parent and sibling
    6.2 Rotate DB parent in DB direction
    6.3 Turn far nephews color to black
    6.4 Remove DB node to single black
*/

// func (d dict[K, V]) Remove(key K) Dict[K, V] {
// 	pt := &d
// 	if d.root == nil {
// 		// Empty tree
// 		return d
// 	} else {
// 		// Find node to delete
// 		dn := d.getNode(key)
// 		removeHelp(pt, dn)
// 		return *pt
// 	}
// }

// func removeHelp[K cmp.Ordered, V any](d *dict[K, V], n *node[K, V]) {
// 	if n == nil {
// 		// Node doesn't exist
// 		return
// 	}
// 	// 2 non-nil children
// 	if n.left != nil && n.right != nil {
// 		suc := findSuccessor(n.right)
// 		n.key = suc.key
// 		n.value = suc.value
// 		// root node
// 		if n.parent == nil {
// 			d.root = n
// 		}
// 		removeHelp(d, suc)
// 		return
// 	}
// 	// 2 nil children
// 	if n.left == nil && n.right == nil {
// 		// root node
// 		if n.parent == nil {
// 			d.root = nil
// 			return
// 		}
//
// 		pSide := parentSide(n, n.parent)
//
// 		switch n.color {
// 		// Case 1 - Red leaf
// 		case RED:
// 			switch pSide {
// 			case LEFT:
// 				// 1.1 - Remove node then exit
// 				n.parent.left = nil
// 				return
// 			case RIGHT:
// 				// 1.1 - Remove node then exit
// 				n.parent.right = nil
// 				return
// 			}
// 		case BLACK:
// 			// Black leaf - DB
// 			fixDB(d, n)
// 			return
// 		}
// 	}
// 	// Black node with red child
// 	if n.left == nil {
// 		// No child on the left
//
// 		// Replace node with child
// 		n.key = n.right.key
// 		n.value = n.right.value
// 		removeHelp(d, n.right)
// 		return
// 	} else {
// 		// No child on the right
//
// 		// Replace node with child
// 		n.key = n.left.key
// 		n.value = n.left.value
// 		removeHelp(d, n.left)
// 		return
// 	}
// }

// func fixDB[K cmp.Ordered, V any](d *dict[K, V], n *node[K, V]) {
// 	// Case 2 - DB is root
// 	if n.parent == nil {
// 		return
// 	}
// 	pColor := n.parent.color
// 	pSide := parentSide(n, n.parent)
// 	sibling := findSibling(n, n.parent)
//
// 	// DB sibling is Black
// 	if sibling.color == BLACK {
// 		// Case 3
// 		if sibling.hasBlackChildren() {
// 			// 3.2 Make sibling red
// 			sibling.color = RED
// 			// Push blackness to parent
// 			n.parent.color = BLACK
// 			if n.hasNilChildren() {
// 				// Node is a leaf node to delete
// 				switch pSide {
// 				case LEFT:
// 					// 3.1 remove node
// 					n.parent.left = nil
// 				case RIGHT:
// 					// 3.1 remove node
// 					n.parent.right = nil
// 				}
// 			}
//
// 			if pColor != BLACK {
// 				return
// 			}
// 			fixDB(d, n.parent)
// 			return
// 		}
// 		switch pSide {
// 		case LEFT:
// 			if sibling.left.isRed() && sibling.right.isBlack() {
// 				// Case 5 - far nephew is Black - near nephew is Red
// 				// 5.1 - Swap colors of sibling and near nephew
// 				sibling.color = RED
// 				sibling.left.color = BLACK
// 				// 5.2 Rotate sibling of DB node in opposite direction of DB node
// 				sibling.srRotation()
// 				// 5.3 Apply Case 6
// 				fixDB(d, n)
// 				return
// 			} else {
// 				// Case 6 - Far nephew is Red
// 				// 6.1 Swap the colors of the DB parent and sibling
// 				n.parent.color = sibling.color
// 				sibling.color = pColor
// 				// 6.2 Rotate DB parent in DB direction
// 				newRoot := n.parent.slRotation()
// 				// 6.3 Turn far nephew's color to black
// 				newRoot.right.color = BLACK
// 				// Check for a new root
// 				if newRoot.parent == nil {
// 					d.root = newRoot
// 				}
//
// 				// 6.4 Remove DB node to single black
// 				if n.hasNilChildren() {
// 					n.parent.left = nil
// 				}
// 				return
// 			}
// 		case RIGHT:
// 			if sibling.right.isRed() && sibling.left.isBlack() {
// 				// Case 5 - far nephew is Black - near nephew is Red
// 				// 5.1 Swap colors of the DB sibling and near nephew
// 				sibling.color = RED
// 				sibling.right.color = BLACK
// 				// 5.2 Rotate sibling of DB node in opposite direction of DB node
// 				sibling.slRotation()
// 				// 5.3 Apply case 6
// 				fixDB(d, n)
// 				return
// 			} else {
// 				// Case 6 - far newphew is Red
// 				// 6.1 Swap the colors of the DB parent and sibling
// 				n.parent.color = sibling.color
// 				sibling.color = pColor
// 				// 6.2 Rotate DB parent in DB direction
// 				newRoot := n.parent.srRotation()
// 				// 6.3 Turn far nephews color to black
// 				newRoot.left.color = BLACK
// 				// Check for a new root
// 				if newRoot.parent == nil {
// 					d.root = newRoot
// 				}
// 				// 6.4 Remove DB node to single black
// 				if n.hasNilChildren() {
// 					n.parent.right = nil
// 				}
// 				return
// 			}
// 		}
// 	}
// 	// Case 4 - Red sibling
//
// 	// 4.1 Swap colors of sibling and parent
// 	sibling.color = pColor
// 	n.parent.color = RED
//
// 	// 4.2 Rotate parent towards n's direction
// 	switch pSide {
// 	case LEFT:
// 		n.parent.slRotation()
// 	case RIGHT:
// 		n.parent.srRotation()
// 	}
// 	if n.hasNilChildren() {
// 		removeHelp(d, n)
// 	}
// 	return
// }

func (n *node[K, V]) hasNilChildren() bool {
	return n.left == nil && n.right == nil
}

func (n *node[K, V]) hasBlackChildren() bool {
	return (n.left == nil || n.left.color == BLACK) && (n.right == nil || n.right.color == BLACK)
}

func (n *node[K, V]) isBlack() bool {
	return n == nil || n.color == BLACK
}
func (n *node[K, V]) isRed() bool {
	return n != nil && n.color == RED
}

func findSuccessor[K cmp.Ordered, V any](n *node[K, V]) *node[K, V] {
	if n.left == nil {
		return n
	}
	return findSuccessor(n.left)
}

// func findSibling[K cmp.Ordered, V any](n *node[K, V], parent *node[K, V]) *node[K, V] {
// 	pDir := parentSide(n, n.parent)
// 	if pDir == LEFT {
// 		return n.parent.right
// 	} else {
// 		return n.parent.left
// 	}
// }

func balance[K cmp.Ordered, V any](n *node[K, V], stk *stack[K, V]) (*node[K, V], *stack[K, V]) {
	// Root case
	if stk.p == nil {
		n.color = BLACK
		return n, stk
	}
	pColor := stk.p.color
	if pColor == BLACK {
		// Nothing more to do
		return n, stk
	}
	// Parent and n are red
	nDir := parentSide(n, stk.p)
	pDir := parentSide(stk.p, stk.pp.p)
	uncle := getUncle(n, stk)
	grandparent := stk.pp.p

	if uncle != nil && uncle.color == RED {
		// Red uncle - push down blackness from grandparent - balance root

		// Copy uncle for mutation
		valU := *uncle
		cpU := &valU
		cpU.color = grandparent.color
		stk.p.color = grandparent.color
		grandparent.color = RED
		setUncle(n, cpU, stk)
		return balance(grandparent, stk.pp.pp)
	}
	// Black uncle
	switch pDir {
	case LEFT:
		switch nDir {
		case LEFT:
			// LL - right rotate on grandparent - balance
			newRoot := stk.pp.p.srRotation(stk.pp.pp)
			rCol := newRoot.right.color
			// Push down newRoot color
			newRoot.right.color = newRoot.color
			newRoot.color = rCol
			// balance newRoot
			return balance(newRoot, stk.pp.pp)
		case RIGHT:
			// LR - rotate parent left - balance left of root
			newRoot := stk.p.slRotation(stk.pp)
			// Add new root to stack as parent
			newStk := &stack[K, V]{p: newRoot, pp: stk.pp}
			return balance(newRoot.left, newStk)
		}
	case RIGHT:
		switch nDir {
		case RIGHT:
			// RR - left rotate on grandparent - balance
			newRoot := stk.pp.p.slRotation(stk.pp.pp)
			// Swap color
			lCol := newRoot.left.color
			newRoot.left.color = newRoot.color
			newRoot.color = lCol
			// balance newRoot
			return balance(newRoot, stk.pp.pp)
		case LEFT:
			//RL - rotate parent right - balance right of root
			newRoot := stk.p.srRotation(stk.pp)
			// Add new root to stack as p
			newStk := &stack[K, V]{p: newRoot, pp: stk.pp}
			return balance(newRoot.right, newStk)
		}
	}
	return n, stk
}

func (x *node[K, V]) srRotation(stk *stack[K, V]) *node[K, V] {
	left := x.left

	// Handle x's parent
	if stk.p != nil {
		pSide := parentSide(x, stk.p)

		switch pSide {
		case LEFT:
			stk.p.left = left
		case RIGHT:
			stk.p.right = left
		}
	}
	// 2. x's left is now lefts right
	x.left = left.right

	// 3. left's right is x
	left.right = x

	return left
}

func (x *node[K, V]) slRotation(stk *stack[K, V]) *node[K, V] {
	right := x.right

	// Handle x's parent
	if stk.p != nil {
		pSide := parentSide(x, stk.p)

		switch pSide {
		case LEFT:
			stk.p.left = right
		case RIGHT:
			stk.p.right = right
		}
	}
	// 2. x's right is right's left
	x.right = right.left

	// 3. right's left is x
	right.left = x

	return right
}

func parentSide[K cmp.Ordered, V any](n *node[K, V], p *node[K, V]) int {
	if p.left != nil && n.key == p.left.key {
		return LEFT
	} else {
		return RIGHT
	}
}

func setUncle[K cmp.Ordered, V any](x *node[K, V], n *node[K, V], stk *stack[K, V]) {
	parent := stk.p
	gp := stk.pp.p

	if parentSide(parent, gp) == LEFT {
		// Uncle is right side
		gp.right = n
	} else {
		gp.left = n
	}
}

func getUncle[K cmp.Ordered, V any](x *node[K, V], stk *stack[K, V]) *node[K, V] {
	parent := stk.p
	gp := stk.pp.p
	if parentSide(parent, gp) == LEFT {
		// Uncle is right side
		return gp.right
	} else {
		return gp.left
	}
}

func getStackRoot[K cmp.Ordered, V any](n *node[K, V], stk *stack[K, V]) *node[K, V] {
	if stk.p == nil {
		return n
	} else {
		return getStackRoot(stk.p, stk.pp)
	}
}

func getStackHelp[K cmp.Ordered, V any](k K, n *node[K, V], st *stack[K, V]) (*stack[K, V], error) {
	switch cmp.Compare(k, n.key) {
	case LT:
		if n.left == nil {
			return st, errors.New("Node does not exist in tree")
		}
		newStk := &stack[K, V]{pp: st, p: n}
		newStk.p = n
		return getStackHelp(k, n.left, newStk)
	case EQ:
		return st, nil
	case GT:
		if n.right == nil {
			return st, errors.New("Node does not exist in tree")
		}
		newStk := &stack[K, V]{pp: st, p: n}
		newStk.p = n
		return getStackHelp(k, n.right, newStk)
	}
	panic("unreachable")
}
