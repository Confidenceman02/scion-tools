// Package dict implements an immutable dictionary, mapping unique keys to values.
// The keys can be any [cmp.Ordered] type.
// Insert, remove, and query operations all take O(log n) time.
package dict

import (
	"cmp"
	"errors"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	. "github.com/Confidenceman02/scion-tools/pkg/maybe"
)

const (
	left = iota
	right
)

const (
	red int = iota
	black
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
type Dict[K, V any] interface {
	rbt() *dict[K, V]
}

/*
Retrieve the internal Red-Black tree
*/
func (d *dict[K, V]) rbt() *dict[K, V] {
	return d
}

type dict[K, V any] struct {
	root *node[K, V]
}

type node[K, V any] struct {
	key   K
	value V
	color int
	left  *node[K, V]
	right *node[K, V]
}

type stack[K, V any] struct {
	pp *stack[K, V]
	p  *node[K, V]
}

type nodeStack[K, V any] struct {
	stack *stack[K, V]
	node  *node[K, V]
}

// BUILD

// Create an empty dictionary.
func Empty[K, V any]() Dict[Comparable[K], V] {
	return &dict[Comparable[K], V]{root: nil}
}

// Create a dictionary with one key-value pair.
func Singleton[K, V any](key Comparable[K], value V) Dict[Comparable[K], V] {
	// Root nodes are always black
	return &dict[Comparable[K], V]{
		root: &node[Comparable[K], V]{
			key:   key,
			value: value,
			color: black,
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
func Insert[K, V any](key Comparable[K], v V, d Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	root := d.rbt().root

	if root == nil {
		return &dict[Comparable[K], V]{root: &node[Comparable[K], V]{key: key, value: v, color: black, left: nil, right: nil}}
	}

	valRoot := *root
	ns := &nodeStack[Comparable[K], V]{node: &valRoot, stack: &stack[Comparable[K], V]{p: nil, pp: nil}}
	insertedNs := insertHelp(key, v, ns)
	newNs := balance(insertedNs)
	rootNs := getNodeStackRoot(newNs)
	return &dict[Comparable[K], V]{root: rootNs.node}
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

// Remove a key-value pair from a dictionary. If the key is not found, no changes are made.
func Remove[K, V any](key Comparable[K], d Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	root := d.rbt().root
	if root == nil {
		// Empty tree
		return d
	}

	// Find nodeStack to delete
	maybeNodeStack := getNodeStack(d.rbt(), key)

	return MaybeWith(
		maybeNodeStack,
		func(j Just[*nodeStack[Comparable[K], V]]) Dict[Comparable[K], V] {
			ns := removeHelp(j.Value)
			rootNs := getNodeStackRoot(ns)
			return &dict[Comparable[K], V]{root: rootNs.node}
		},
		func(n Nothing) Dict[Comparable[K], V] { return d },
	)
}

// Get a Just nodeStack or Nothing if node doesn't exist
func getNodeStack[K, V any](d *dict[Comparable[K], V], targetKey Comparable[K]) Maybe[*nodeStack[Comparable[K], V]] {
	if d.root == nil {
		return Nothing{}
	} else {
		valRoot := *d.root
		return getNodeStackHelp(targetKey, &nodeStack[Comparable[K], V]{stack: &stack[Comparable[K], V]{p: nil, pp: nil}, node: &valRoot})
	}
}

/*
Gets a 'Just' nodeStack or 'Nothing' if it doesn't exist
*/
func getNodeStackHelp[K, V any](targetKey Comparable[K], ns *nodeStack[Comparable[K], V]) Maybe[*node[Comparable[K], V]] {
getNodeStackHelpL:
	for {
		switch targetKey.Cmp(ns.node.key) {
		case -1:
			if ns.node.left == nil {
				return Nothing{}
			} else {
				newStack := &stack[Comparable[K], V]{pp: ns.stack, p: ns.node}
				valL := *ns.node.left
				ns.node.left = &valL
				newNs := &nodeStack[Comparable[K], V]{node: ns.node.left, stack: newStack}
				ns = newNs
				continue getNodeStackHelpL
			}
		case 0:
			return Just[*nodeStack[Comparable[K], V]]{Value: ns}
		case +1:
			if ns.node.right == nil {
				return Nothing{}
			} else {
				newStack := &stack[Comparable[K], V]{pp: ns.stack, p: ns.node}
				valR := *ns.node.right
				ns.node.right = &valR
				newNs := &nodeStack[Comparable[K], V]{node: ns.node.right, stack: newStack}
				ns = newNs
				continue getNodeStackHelpL
			}
		}
	}
}

func insertHelp[K, V any](key Comparable[K], value V, ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
insertHelpL:
	for {
		nKey := ns.node.key
		switch key.Cmp(nKey) {
		case -1:
			if ns.node.left == nil {
				ns.node.left = &node[Comparable[K], V]{key: key, value: value, color: red, left: nil, right: nil}
				newStk := &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}
				newNs := &nodeStack[Comparable[K], V]{node: ns.node.left, stack: newStk}
				return newNs
			} else {
				valL := *ns.node.left
				ns.node.left = &valL
				newStk := &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}
				tempNs := &nodeStack[Comparable[K], V]{node: ns.node.left, stack: newStk}
				ns = tempNs
				continue insertHelpL
			}
		case 0:
			ns.node.value = value
			return ns
		case +1:
			if ns.node.right == nil {
				ns.node.right = &node[Comparable[K], V]{key: key, value: value, color: red, left: nil, right: nil}
				newStk := &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}
				newNs := &nodeStack[Comparable[K], V]{node: ns.node.right, stack: newStk}
				return newNs
			} else {
				valR := *ns.node.right
				ns.node.right = &valR
				newStk := &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}
				newNs := &nodeStack[Comparable[K], V]{node: ns.node.right, stack: newStk}
				ns = newNs
				continue insertHelpL
			}
		}
	}
}

func removeHelp[K, V any](ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
removeHelpL:
	for {
		// 2 non-nil children
		if ns.node.left != nil && ns.node.right != nil {
			// Copy right node
			valR := *ns.node.right
			ns.node.right = &valR

			// Create new stack for right child
			rightStack := &nodeStack[Comparable[K], V]{node: ns.node.right, stack: &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}}

			// Find in order successor
			succStk := findSuccessor(rightStack)

			// Swap key and values
			ns.node.key = succStk.node.key
			ns.node.value = succStk.node.value

			ns = succStk
			// Find new case
			continue removeHelpL
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
			case red:
				switch pSide {
				case left:
					// 1.1 - remove node then exit
					ns.stack.p.left = nil
					return ns
				case right:
					// 1.1 - remove node then exit
					ns.stack.p.right = nil
					return ns
				}
			case black:
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
			newStack := &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}
			newNs := &nodeStack[Comparable[K], V]{node: ns.node.right, stack: newStack}
			ns = newNs
			continue removeHelpL
		} else {
			// Copy left node
			valL := *ns.node.left
			ns.node.left = &valL

			// Replace node with right node
			ns.node.key = ns.node.left.key
			ns.node.value = ns.node.left.value

			// New nodeStack
			newStack := &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}
			newNs := &nodeStack[Comparable[K], V]{node: ns.node.left, stack: newStack}
			ns = newNs
			continue removeHelpL
		}
	}
}

func fixDB[K, V any](ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
fixDBL:
	for {
		// Case 2 - DB is root
		if ns.stack.p == nil {
			return ns
		}
		pColor := ns.stack.p.color
		pSide := parentSide(ns)
		sNs := findSibling(ns)

		// DB sibling is Black
		if sNs.node.color == black {

			// Case 3
			if sNs.node.hasBlackChildren() {
				// 3.1 Remove node if leaf
				if ns.node.hasNilChildren() {
					switch pSide {
					case right:
						ns.stack.p.right = nil
					case left:
						ns.stack.p.left = nil
					}
				}
				// 3.2 Make sibling red
				sNs.node.color = red
				// 3.3 Push blackness to parent
				sNs.stack.p.color = black
				// Check if parent is DB
				if pColor != black {
					return ns
				} else {
					ns = &nodeStack[Comparable[K], V]{stack: ns.stack.pp, node: ns.stack.p}
					continue fixDBL
				}
			}
			switch pSide {
			case left:
				if sNs.node.left.isRed() && sNs.node.right.isBlack() {
					// Case 5 - far nephew is black - near nephew is red
					// Copy near nephew
					valSL := *sNs.node.left
					// 5.1 - Swap colors of sibling and near nephew
					valSL.color = black
					sNs.node.left = &valSL
					sNs.node.color = red

					// 5.2 Rotate sibling of DB node in opposite direction of DB node
					srRotationV2(sNs.node, sNs.stack)
					// 5.3 Apply Case 6
					continue fixDBL
				} else {
					// Case 6 - Far nephew is Red
					valSR := *sNs.node.right
					// 6.1 Swap the colors of the DB parent and sibling
					ns.stack.p.color = sNs.node.color
					sNs.node.color = pColor
					// 6.2 Rotate DB parent in DB direction
					grandparentStk := ns.stack.pp
					newRoot := slRotationV2(ns.stack.p, grandparentStk)
					newNs := &nodeStack[Comparable[K], V]{stack: grandparentStk, node: newRoot}
					// 6.3 Turn far nephew's color to black
					valSR.color = black
					newNs.node.right = &valSR
					// 6.4 Remove DB node to single black
					if ns.node.hasNilChildren() {
						ns.stack.p.left = nil
					}
					return newNs
				}
			case right:
				if sNs.node.left.isBlack() && sNs.node.right.isRed() {
					// Case 5 - far nephew is black - near nephew is red
					// Copy near nephew
					valSR := *sNs.node.right
					// 5.1 - Swap colors of sibling and near nephew
					valSR.color = black
					sNs.node.right = &valSR
					sNs.node.color = red

					// 5.2 Rotate sibling of DB node in opposite direction of DB node
					slRotationV2(sNs.node, sNs.stack)
					// 5.3 Apply Case 6
					continue fixDBL
				} else {
					// Case 6 - Far nephew is Red
					valSL := *sNs.node.left
					// 6.1 Swap the colors of the DB parent and sibling
					ns.stack.p.color = sNs.node.color
					sNs.node.color = pColor
					// 6.2 Rotate DB parent in DB direction
					grandparentStk := ns.stack.pp
					newRoot := srRotationV2(ns.stack.p, grandparentStk)
					newNs := &nodeStack[Comparable[K], V]{stack: grandparentStk, node: newRoot}
					// 6.3 Turn far nephew's color to black
					valSL.color = black
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
		var newRoot = &node[Comparable[K], V]{}
		switch pSide {
		case left:
			newRoot = slRotationV2(ns.stack.p, grandparentStk)

		case right:
			newRoot = srRotationV2(ns.stack.p, grandparentStk)
		}
		// Add new root to stack as parent
		newRootStk := &stack[Comparable[K], V]{p: newRoot, pp: ns.stack.pp}
		newPStk := &stack[Comparable[K], V]{p: ns.stack.p, pp: newRootStk}
		ns.stack = newPStk
		continue fixDBL
	}
}

func (n *node[K, V]) hasNilChildren() bool {
	return n.left == nil && n.right == nil
}

func (n *node[K, V]) hasBlackChildren() bool {
	return (n.left == nil || n.left.color == black) && (n.right == nil || n.right.color == black)
}

func (n *node[K, V]) isBlack() bool {
	return n == nil || n.color == black
}
func (n *node[K, V]) isRed() bool {
	return n != nil && n.color == red
}

func findSuccessor[K, V any](ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
	if ns.node.left == nil {
		return ns
	} else {
		valL := *ns.node.left
		ns.node.left = &valL
		newStack := &nodeStack[Comparable[K], V]{node: ns.node.left, stack: &stack[Comparable[K], V]{p: ns.node, pp: ns.stack}}
		return findSuccessor(newStack)
	}
}

// Find sibling nodeStack
func findSibling[K, V any](ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
	pDir := parentSide(ns)
	if pDir == left {
		valR := *ns.stack.p.right
		ns.stack.p.right = &valR
		return &nodeStack[Comparable[K], V]{stack: ns.stack, node: ns.stack.p.right}
	} else {
		valL := *ns.stack.p.left
		ns.stack.p.left = &valL
		return &nodeStack[Comparable[K], V]{stack: ns.stack, node: ns.stack.p.left}
	}
}

func balance[K, V any](ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
balanceL:
	for {
		// Root case
		if ns.stack.p == nil {
			ns.node.color = black
			return ns
		}
		pColor := ns.stack.p.color
		if pColor == black {
			// Nothing more to do
			return ns
		}
		// Parent and node are red
		nDir := parentSide(ns)
		pDir := parentSide(&nodeStack[Comparable[K], V]{node: ns.stack.p, stack: ns.stack.pp})
		uncle := getUncle(ns)
		grandparent := ns.stack.pp.p

		if uncle != nil && uncle.color == red {
			// Red uncle - push down blackness from grandparent - balance root

			// Copy uncle for mutation
			valU := *uncle
			cpU := &valU

			cpU.color = grandparent.color
			ns.stack.p.color = grandparent.color
			grandparent.color = red
			setUncle(ns, cpU)
			newNs := &nodeStack[Comparable[K], V]{node: grandparent, stack: ns.stack.pp.pp}
			ns = newNs
			continue balanceL
		}
		// Black uncle
		switch pDir {
		case left:
			switch nDir {
			case left:
				// LL - right rotate on grandparent - balance
				newRoot := srRotationV2(ns.stack.pp.p, ns.stack.pp.pp)

				// Swap colors
				rCol := newRoot.right.color
				newRoot.right.color = newRoot.color
				newRoot.color = rCol

				// balance newRoot
				newNs := &nodeStack[Comparable[K], V]{node: newRoot, stack: ns.stack.pp.pp}
				ns = newNs
				continue balanceL
			case right:
				// LR - rotate parent left - balance left of root
				newRoot := slRotationV2(ns.stack.p, ns.stack.pp)

				// Add new root to stack as parent
				newStk := &stack[Comparable[K], V]{p: newRoot, pp: ns.stack.pp}
				newNs := &nodeStack[Comparable[K], V]{node: newRoot.left, stack: newStk}
				ns = newNs
				continue balanceL
			}
		case right:
			switch nDir {
			case right:
				// RR - left rotate on grandparent - balance
				newRoot := slRotationV2(ns.stack.pp.p, ns.stack.pp.pp)

				// Swap color
				lCol := newRoot.left.color
				newRoot.left.color = newRoot.color
				newRoot.color = lCol

				// balance newRoot
				newNs := &nodeStack[Comparable[K], V]{node: newRoot, stack: ns.stack.pp.pp}
				ns = newNs
				continue balanceL
			case left:
				//RL - rotate parent right - balance right of root
				newRoot := srRotationV2(ns.stack.p, ns.stack.pp)

				// Add newRoot to stack as parent
				newStk := &stack[Comparable[K], V]{p: newRoot, pp: ns.stack.pp}
				newNs := &nodeStack[Comparable[K], V]{node: newRoot.right, stack: newStk}
				ns = newNs
				continue balanceL
			}
		}
		return ns
	}
}

func srRotationV2[K, V any](x *node[Comparable[K], V], stk *stack[Comparable[K], V]) *node[Comparable[K], V] {
	leftNode := x.left

	// Handle x's parent
	if stk.p != nil {
		pSide := parentSide(&nodeStack[Comparable[K], V]{node: x, stack: stk})

		switch pSide {
		case left:
			stk.p.left = leftNode
		case right:
			stk.p.right = leftNode
		}
	}
	// 2. x's left is now lefts right
	x.left = leftNode.right

	// 3. left's right is x
	leftNode.right = x

	return leftNode
}

func slRotationV2[K, V any](x *node[Comparable[K], V], stk *stack[Comparable[K], V]) *node[Comparable[K], V] {
	rightNode := x.right

	// Handle x's parent
	if stk.p != nil {
		pSide := parentSide(&nodeStack[Comparable[K], V]{node: x, stack: stk})

		switch pSide {
		case left:
			stk.p.left = rightNode
		case right:
			stk.p.right = rightNode
		}
	}
	// 2. x's right is right's left
	x.right = rightNode.left

	// 3. right's left is x
	rightNode.left = x

	return rightNode
}

func parentSide[K, V any](ns *nodeStack[Comparable[K], V]) int {
	p := ns.stack.p
	if p.left != nil && ns.node.key == p.left.key {
		return left
	} else {
		return right
	}
}

func setUncle[K, V any](ns *nodeStack[Comparable[K], V], unc *node[Comparable[K], V]) {
	parent := ns.stack.p
	gp := ns.stack.pp.p

	if parentSide(&nodeStack[Comparable[K], V]{node: parent, stack: ns.stack.pp}) == left {
		// Uncle is right side
		gp.right = unc
	} else {
		gp.left = unc
	}
}

func getUncle[K, V any](ns *nodeStack[Comparable[K], V]) *node[Comparable[K], V] {
	parent := ns.stack.p
	gp := ns.stack.pp.p
	if parentSide(&nodeStack[Comparable[K], V]{node: parent, stack: ns.stack.pp}) == left {
		// Uncle is right side
		return gp.right
	} else {
		return gp.left
	}
}

func getNodeStackRoot[K, V any](ns *nodeStack[Comparable[K], V]) *nodeStack[Comparable[K], V] {
	if ns.stack.p == nil {
		return ns
	} else {
		newNs := &nodeStack[Comparable[K], V]{node: ns.stack.p, stack: ns.stack.pp}
		return getNodeStackRoot(newNs)
	}
}

func getStackHelp[K cmp.Ordered, V any](k K, n *node[K, V], st *stack[K, V]) (*stack[K, V], error) {
getNodeStackHelpL:
	for {
		switch cmp.Compare(k, n.key) {
		case -1:
			if n.left == nil {
				return st, errors.New("Node does not exist in tree")
			}
			newStk, newN := &stack[K, V]{pp: st, p: n}, n.left
			n = newN
			st = newStk
			continue getNodeStackHelpL
		case 0:
			return st, nil
		case +1:
			if n.right == nil {
				return st, errors.New("Node does not exist in tree")
			}
			newStk, newN := &stack[K, V]{pp: st, p: n}, n.right
			n = newN
			st = newStk
			continue getNodeStackHelpL
		}
	}
}

// QUERY

// Determine if a dictionary is empty.
func IsEmpty[K, V any](d Dict[Comparable[K], V]) bool {
	return d.rbt().root == nil
}

// Member - Determine if a key is in a dictionary.
func Member[K, V any](k Comparable[K], d Dict[Comparable[K], V]) bool {
	root := d.rbt().root

	if root == nil {
		return false
	} else {
		return memberHelp(k, root)
	}
}

func memberHelp[K, V any](k Comparable[K], n *node[Comparable[K], V]) bool {
memberHelpL:
	for {
		if n == nil {
			return false
		} else {
			switch k.Cmp(n.key) {
			case -1:
				n = n.left
				continue memberHelpL
			case 0:
				return true
			case +1:
				n = n.right
				continue memberHelpL
			}
		}
	}
}

// Get the value associated with a key.
// If the key is not found, return [Nothing].
// This is useful when you are not sure if a key will be in the dictionary.
func Get[K, V any](targetKey Comparable[K], d Dict[Comparable[K], V]) Maybe[V] {
	root := d.rbt().root
	if root == nil {
		return Nothing{}
	} else {
		return getHelp(targetKey, root)
	}
}

func getHelp[K, V any](targetKey Comparable[K], n *node[Comparable[K], V]) Maybe[V] {
getHelpL:
	for {
		if n == nil {
			return Nothing{}
		} else {
			switch targetKey.Cmp(n.key) {
			case -1:
				n = n.left
				continue getHelpL
			case 0:
				return Just[V]{Value: n.value}
			case +1:
				n = n.right
				continue getHelpL
			}
		}
	}
}

// LISTS
// TRANSFORM
// COMBINE
