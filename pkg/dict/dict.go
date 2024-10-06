// Package dict implements an immutable dictionary, mapping unique keys to values.
// The keys can be any [cmp.Ordered] type.
// Insert, remove, and query operations all take O(log n) time.
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

// Dict represents a dictionary of keys and values.
// So a Dict[string, User] is a dictionary that lets you look up a string (such as user names) and find the associated User.
type Dict[K cmp.Ordered, V any] interface {
	rbt() *dict[K, V]
	getNodeStack(k K) maybe.Maybe[*nodeStack[K, V]]
}

/*
Retrieve the internal Red-Black tree
*/
func (d *dict[K, V]) rbt() *dict[K, V] {
	return d
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

type nodeStack[K cmp.Ordered, V any] struct {
	stack *stack[K, V]
	node  *node[K, V]
}

// BUILDERS

// Create an empty dictionary.
func Empty[K cmp.Ordered, V any]() Dict[K, V] {
	return &dict[K, V]{root: nil}
}

// Create a dictionary with one key-value pair.
func Singleton[K cmp.Ordered, V any](key K, value V) Dict[K, V] {
	// Root nodes are always black
	return &dict[K, V]{
		root: &node[K, V]{
			key:   key,
			value: value,
			color: BLACK,
			left:  nil,
			right: nil},
	}
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

// Insert a key-value pair into a dictionary. Replaces the value when there is a collision.
func Insert[K cmp.Ordered, V any](key K, v V, d Dict[K, V]) Dict[K, V] {
	root := d.rbt().root

	if root == nil {
		return &dict[K, V]{root: &node[K, V]{key: key, value: v, color: BLACK, left: nil, right: nil}}
	}

	valRoot := *root
	ns := &nodeStack[K, V]{node: &valRoot, stack: &stack[K, V]{p: nil, pp: nil}}
	insertedNs := insertHelp(key, v, ns)
	newNs := balance(insertedNs)
	rootNs := getNodeStackRoot(newNs)
	return &dict[K, V]{root: rootNs.node}
}

// TODO Update

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

// Remove a key-value pair from a dictionary. If the key is not found, no changes are made.
func Remove[K cmp.Ordered, V any](key K, d Dict[K, V]) Dict[K, V] {
	root := d.rbt().root
	if root == nil {
		// Empty tree
		return d
	}

	// Find nodeStack to delete
	maybeNodeStack := d.getNodeStack(key)

	return maybe.MaybeWith(
		maybeNodeStack,
		func(j maybe.Just[*nodeStack[K, V]]) Dict[K, V] {
			ns := removeHelp(j.Value)
			rootNs := getNodeStackRoot(ns)
			return &dict[K, V]{root: rootNs.node}
		},
		func(n maybe.Nothing) Dict[K, V] { return d },
	)
}

// QUERY

// Determine if a dictionary is empty.
func IsEmpty[K cmp.Ordered, V any](d Dict[K, V]) bool {
	return d.rbt().root == nil
}

// Member - Determine if a key is in a dictionary.
func Member[K cmp.Ordered, V any](k K, d Dict[K, V]) bool {
	root := d.rbt().root

	if root == nil {
		return false
	} else {
		return memberHelp(k, root)
	}
}

func memberHelp[K cmp.Ordered, V any](k K, n *node[K, V]) bool {
	if n != nil {
		switch cmp.Compare(k, n.key) {
		case LT:
			return memberHelp(k, n.left)
		case EQ:
			return true
		case GT:
			return memberHelp(k, n.right)
		}
	}
	return false
}

// Get the value associated with a key.
// If the key is not found, return [maybe.Nothing].
// This is useful when you are not sure if a key will be in the dictionary.
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

// Get a Just nodeStack or Nothing if node doesn't exist
func (d *dict[K, V]) getNodeStack(targetKey K) maybe.Maybe[*nodeStack[K, V]] {
	if d.root == nil {
		return maybe.Nothing{}
	} else {
		valRoot := *d.root
		return getNodeStackHelp(targetKey, &nodeStack[K, V]{stack: &stack[K, V]{p: nil, pp: nil}, node: &valRoot})
	}
}

/*
Gets a 'Just' nodeStack or 'Nothing' if it doesn't exist
*/
func getNodeStackHelp[K cmp.Ordered, V any](targetKey K, ns *nodeStack[K, V]) maybe.Maybe[*node[K, V]] {
	switch cmp.Compare(targetKey, ns.node.key) {
	case LT:
		if ns.node.left == nil {
			return maybe.Nothing{}
		} else {
			newStack := &stack[K, V]{pp: ns.stack, p: ns.node}
			valL := *ns.node.left
			ns.node.left = &valL
			newNs := &nodeStack[K, V]{node: ns.node.left, stack: newStack}
			return getNodeStackHelp(targetKey, newNs)
		}
	case EQ:
		return maybe.Just[*nodeStack[K, V]]{Value: ns}
	case GT:
		if ns.node.right == nil {
			return maybe.Nothing{}
		} else {
			newStack := &stack[K, V]{pp: ns.stack, p: ns.node}
			valR := *ns.node.right
			ns.node.right = &valR
			newNs := &nodeStack[K, V]{node: ns.node.right, stack: newStack}
			return getNodeStackHelp(targetKey, newNs)
		}
	}
	panic("getNodeHelp unreachable")
}

func insertHelp[K cmp.Ordered, V any](key K, value V, ns *nodeStack[K, V]) *nodeStack[K, V] {
	nKey := ns.node.key
	switch cmp.Compare(key, nKey) {
	case LT:
		if ns.node.left == nil {
			ns.node.left = &node[K, V]{key: key, value: value, color: RED, left: nil, right: nil}
			newStk := &stack[K, V]{p: ns.node, pp: ns.stack}
			newNs := &nodeStack[K, V]{node: ns.node.left, stack: newStk}
			return newNs
		} else {
			valL := *ns.node.left
			ns.node.left = &valL
			newStk := &stack[K, V]{p: ns.node, pp: ns.stack}
			newNs := &nodeStack[K, V]{node: ns.node.left, stack: newStk}
			return insertHelp(key, value, newNs)
		}
	case EQ:
		ns.node.value = value
		return ns
	case GT:
		if ns.node.right == nil {
			ns.node.right = &node[K, V]{key: key, value: value, color: RED, left: nil, right: nil}
			newStk := &stack[K, V]{p: ns.node, pp: ns.stack}
			newNs := &nodeStack[K, V]{node: ns.node.right, stack: newStk}
			return newNs
		} else {
			valR := *ns.node.right
			ns.node.right = &valR
			newStk := &stack[K, V]{p: ns.node, pp: ns.stack}
			newNs := &nodeStack[K, V]{node: ns.node.right, stack: newStk}
			return insertHelp(key, value, newNs)
		}
	}
	panic("unreachable")
}

func removeHelp[K cmp.Ordered, V any](ns *nodeStack[K, V]) *nodeStack[K, V] {
	// 2 non-nil children
	if ns.node.left != nil && ns.node.right != nil {
		// Copy right node
		valR := *ns.node.right
		ns.node.right = &valR

		// Create new stack for right child
		rightStack := &nodeStack[K, V]{node: ns.node.right, stack: &stack[K, V]{p: ns.node, pp: ns.stack}}

		// Find in order successor
		succStk := findSuccessor(rightStack)

		// Swap key and values
		ns.node.key = succStk.node.key
		ns.node.value = succStk.node.value

		// Find new case
		return removeHelp(succStk)
	}
	// 2 nil children
	if ns.node.left == nil && ns.node.right == nil {
		// root node
		if ns.stack.p == nil {
			ns.node = nil
			return ns
		}
		pSide := parentSide(ns)

		switch ns.node.color {
		// case 1 - red leaf
		case RED:
			switch pSide {
			case LEFT:
				// 1.1 - remove node then exit
				ns.stack.p.left = nil
				return ns
			case RIGHT:
				// 1.1 - remove node then exit
				ns.stack.p.right = nil
				return ns
			}
		case BLACK:
			return fixDB(ns)
		}
	}
	// Black node with red child
	if ns.node.left == nil {
		// Copy right node
		valR := *ns.node.right
		ns.node.right = &valR

		// Replace node with right node
		ns.node.key = ns.node.right.key
		ns.node.value = ns.node.right.value

		// New nodeStack
		newStack := &stack[K, V]{p: ns.node, pp: ns.stack}
		newNs := &nodeStack[K, V]{node: ns.node.right, stack: newStack}

		return removeHelp(newNs)
	} else {
		// Copy left node
		valL := *ns.node.left
		ns.node.left = &valL

		// Replace node with right node
		ns.node.key = ns.node.left.key
		ns.node.value = ns.node.left.value

		// New nodeStack
		newStack := &stack[K, V]{p: ns.node, pp: ns.stack}
		newNs := &nodeStack[K, V]{node: ns.node.left, stack: newStack}

		return removeHelp(newNs)
	}
}

func fixDB[K cmp.Ordered, V any](ns *nodeStack[K, V]) *nodeStack[K, V] {
	// Case 2 - DB is root
	if ns.stack.p == nil {
		return ns
	}
	pColor := ns.stack.p.color
	pSide := parentSide(ns)
	sNs := findSibling(ns)

	// DB sibling is Black
	if sNs.node.color == BLACK {

		// Case 3
		if sNs.node.hasBlackChildren() {
			// 3.1 Remove node if leaf
			if ns.node.hasNilChildren() {
				switch pSide {
				case RIGHT:
					ns.stack.p.right = nil
				case LEFT:
					ns.stack.p.left = nil
				}
			}
			// 3.2 Make sibling red
			sNs.node.color = RED
			// 3.3 Push blackness to parent
			sNs.stack.p.color = BLACK
			// Check if parent is DB
			if pColor != BLACK {
				return ns
			} else {
				return fixDB(&nodeStack[K, V]{stack: ns.stack.pp, node: ns.stack.p})
			}
		}
		switch pSide {
		case LEFT:
			if sNs.node.left.isRed() && sNs.node.right.isBlack() {
				// Case 5 - far nephew is BLACK - near nephew is RED
				// Copy near nephew
				valSL := *sNs.node.left
				// 5.1 - Swap colors of sibling and near nephew
				valSL.color = BLACK
				sNs.node.left = &valSL
				sNs.node.color = RED

				// 5.2 Rotate sibling of DB node in opposite direction of DB node
				sNs.node.srRotation(sNs.stack)
				// 5.3 Apply Case 6
				return fixDB(ns)
			} else {
				// Case 6 - Far nephew is Red
				valSR := *sNs.node.right
				// 6.1 Swap the colors of the DB parent and sibling
				ns.stack.p.color = sNs.node.color
				sNs.node.color = pColor
				// 6.2 Rotate DB parent in DB direction
				grandparentStk := ns.stack.pp
				newRoot := ns.stack.p.slRotation(grandparentStk)
				newNs := &nodeStack[K, V]{stack: grandparentStk, node: newRoot}
				// 6.3 Turn far nephew's color to black
				valSR.color = BLACK
				newNs.node.right = &valSR
				// 6.4 Remove DB node to single black
				if ns.node.hasNilChildren() {
					ns.stack.p.left = nil
				}
				return newNs
			}
		case RIGHT:
			if sNs.node.left.isBlack() && sNs.node.right.isRed() {
				// Case 5 - far nephew is BLACK - near nephew is RED
				// Copy near nephew
				valSR := *sNs.node.right
				// 5.1 - Swap colors of sibling and near nephew
				valSR.color = BLACK
				sNs.node.right = &valSR
				sNs.node.color = RED

				// 5.2 Rotate sibling of DB node in opposite direction of DB node
				sNs.node.slRotation(sNs.stack)
				// 5.3 Apply Case 6
				return fixDB(ns)
			} else {
				// Case 6 - Far nephew is Red
				valSL := *sNs.node.left
				// 6.1 Swap the colors of the DB parent and sibling
				ns.stack.p.color = sNs.node.color
				sNs.node.color = pColor
				// 6.2 Rotate DB parent in DB direction
				grandparentStk := ns.stack.pp
				newRoot := ns.stack.p.srRotation(grandparentStk)
				newNs := &nodeStack[K, V]{stack: grandparentStk, node: newRoot}
				// 6.3 Turn far nephew's color to black
				valSL.color = BLACK
				newNs.node.left = &valSL
				// 6.4 Remove DB node to single black
				if ns.node.hasNilChildren() {
					ns.stack.p.right = nil
				}
				return newNs
			}
		}
	}
	// Case 4 - Red sibling
	// 4.1 Swap colors of sibling and parent
	ns.stack.p.color = sNs.node.color
	sNs.node.color = pColor

	// 4.2 Rotate parent towards n's direction
	grandparentStk := ns.stack.pp
	var newRoot = &node[K, V]{}
	switch pSide {
	case LEFT:
		newRoot = ns.stack.p.slRotation(grandparentStk)

	case RIGHT:
		newRoot = ns.stack.p.srRotation(grandparentStk)
	}
	// Add new root to stack as parent
	newRootStk := &stack[K, V]{p: newRoot, pp: ns.stack.pp}
	newPStk := &stack[K, V]{p: ns.stack.p, pp: newRootStk}
	ns.stack = newPStk
	return fixDB(ns)
}

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

func findSuccessor[K cmp.Ordered, V any](ns *nodeStack[K, V]) *nodeStack[K, V] {
	if ns.node.left == nil {
		return ns
	} else {
		valL := *ns.node.left
		ns.node.left = &valL
		newStack := &nodeStack[K, V]{node: ns.node.left, stack: &stack[K, V]{p: ns.node, pp: ns.stack}}
		return findSuccessor(newStack)
	}
}

// Find sibling nodeStack
func findSibling[K cmp.Ordered, V any](ns *nodeStack[K, V]) *nodeStack[K, V] {
	pDir := parentSide(ns)
	if pDir == LEFT {
		valR := *ns.stack.p.right
		ns.stack.p.right = &valR
		return &nodeStack[K, V]{stack: ns.stack, node: ns.stack.p.right}
	} else {
		valL := *ns.stack.p.left
		ns.stack.p.left = &valL
		return &nodeStack[K, V]{stack: ns.stack, node: ns.stack.p.left}
	}
}

func balance[K cmp.Ordered, V any](ns *nodeStack[K, V]) *nodeStack[K, V] {
	// Root case
	if ns.stack.p == nil {
		ns.node.color = BLACK
		return ns
	}
	pColor := ns.stack.p.color
	if pColor == BLACK {
		// Nothing more to do
		return ns
	}
	// Parent and node are red
	nDir := parentSide(ns)
	pDir := parentSide(&nodeStack[K, V]{node: ns.stack.p, stack: ns.stack.pp})
	uncle := getUncle(ns)
	grandparent := ns.stack.pp.p

	if uncle != nil && uncle.color == RED {
		// Red uncle - push down blackness from grandparent - balance root

		// Copy uncle for mutation
		valU := *uncle
		cpU := &valU

		cpU.color = grandparent.color
		ns.stack.p.color = grandparent.color
		grandparent.color = RED
		setUncle(ns, cpU)
		newNs := &nodeStack[K, V]{node: grandparent, stack: ns.stack.pp.pp}
		return balance(newNs)
	}
	// Black uncle
	switch pDir {
	case LEFT:
		switch nDir {
		case LEFT:
			// LL - right rotate on grandparent - balance
			newRoot := ns.stack.pp.p.srRotation(ns.stack.pp.pp)

			// Swap colors
			rCol := newRoot.right.color
			newRoot.right.color = newRoot.color
			newRoot.color = rCol

			// balance newRoot
			newNs := &nodeStack[K, V]{node: newRoot, stack: ns.stack.pp.pp}
			return balance(newNs)
		case RIGHT:
			// LR - rotate parent left - balance left of root
			newRoot := ns.stack.p.slRotation(ns.stack.pp)

			// Add new root to stack as parent
			newStk := &stack[K, V]{p: newRoot, pp: ns.stack.pp}
			newNs := &nodeStack[K, V]{node: newRoot.left, stack: newStk}
			return balance(newNs)
		}
	case RIGHT:
		switch nDir {
		case RIGHT:
			// RR - left rotate on grandparent - balance
			newRoot := ns.stack.pp.p.slRotation(ns.stack.pp.pp)

			// Swap color
			lCol := newRoot.left.color
			newRoot.left.color = newRoot.color
			newRoot.color = lCol

			// balance newRoot
			newNs := &nodeStack[K, V]{node: newRoot, stack: ns.stack.pp.pp}
			return balance(newNs)
		case LEFT:
			//RL - rotate parent right - balance right of root
			newRoot := ns.stack.p.srRotation(ns.stack.pp)

			// Add newRoot to stack as parent
			newStk := &stack[K, V]{p: newRoot, pp: ns.stack.pp}
			newNs := &nodeStack[K, V]{node: newRoot.right, stack: newStk}
			return balance(newNs)
		}
	}
	return ns
}

func (x *node[K, V]) srRotation(stk *stack[K, V]) *node[K, V] {
	left := x.left

	// Handle x's parent
	if stk.p != nil {
		pSide := parentSide(&nodeStack[K, V]{node: x, stack: stk})

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
		pSide := parentSide(&nodeStack[K, V]{node: x, stack: stk})

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

func parentSide[K cmp.Ordered, V any](ns *nodeStack[K, V]) int {
	p := ns.stack.p
	if p.left != nil && ns.node.key == p.left.key {
		return LEFT
	} else {
		return RIGHT
	}
}

func setUncle[K cmp.Ordered, V any](ns *nodeStack[K, V], unc *node[K, V]) {
	parent := ns.stack.p
	gp := ns.stack.pp.p

	if parentSide(&nodeStack[K, V]{node: parent, stack: ns.stack.pp}) == LEFT {
		// Uncle is right side
		gp.right = unc
	} else {
		gp.left = unc
	}
}

func getUncle[K cmp.Ordered, V any](ns *nodeStack[K, V]) *node[K, V] {
	parent := ns.stack.p
	gp := ns.stack.pp.p
	if parentSide(&nodeStack[K, V]{node: parent, stack: ns.stack.pp}) == LEFT {
		// Uncle is right side
		return gp.right
	} else {
		return gp.left
	}
}

func getNodeStackRoot[K cmp.Ordered, V any](ns *nodeStack[K, V]) *nodeStack[K, V] {
	if ns.stack.p == nil {
		return ns
	} else {
		newNs := &nodeStack[K, V]{node: ns.stack.p, stack: ns.stack.pp}
		return getNodeStackRoot(newNs)
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
