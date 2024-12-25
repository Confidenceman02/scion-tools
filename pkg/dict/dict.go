// Package dict implements an immutable dictionary, mapping unique keys to values.
// The keys can be any [cmp.Ordered] type.
// Insert, remove, and query operations all take O(log n) time.
package dict

import (
	"cmp"
	"errors"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
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

// Update the value of a dictionary for a specific key with a given function.
func Update[K, V any](targetKey Comparable[K], f func(maybe.Maybe[V]) maybe.Maybe[V], d Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	return maybe.MaybeWith(
		f(Get(targetKey, d)),
		func(j maybe.Just[V]) Dict[Comparable[K], V] { return Insert(targetKey, j.Value, d) },
		func(n maybe.Nothing) Dict[Comparable[K], V] { return Remove(targetKey, d) },
	)
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
	} else {
		valRoot := *root
		ns := &nodeStack[Comparable[K], V]{node: &valRoot, stack: &stack[Comparable[K], V]{p: nil, pp: nil}}
		insertedNs := insertHelp(key, v, ns)
		newNs := balance(insertedNs)
		rootNs := getNodeStackRoot(newNs)
		return &dict[Comparable[K], V]{root: rootNs.node}
	}
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

	return maybe.MaybeWith(
		maybeNodeStack,
		func(j maybe.Just[*nodeStack[Comparable[K], V]]) Dict[Comparable[K], V] {
			ns := removeHelp(j.Value)
			rootNs := getNodeStackRoot(ns)
			return &dict[Comparable[K], V]{root: rootNs.node}
		},
		func(n maybe.Nothing) Dict[Comparable[K], V] { return d },
	)
}

// Get a Just nodeStack or Nothing if node doesn't exist
func getNodeStack[K, V any](d *dict[Comparable[K], V], targetKey Comparable[K]) maybe.Maybe[*nodeStack[Comparable[K], V]] {
	if d.root == nil {
		return maybe.Nothing{}
	} else {
		valRoot := *d.root
		return getNodeStackHelp(targetKey, &nodeStack[Comparable[K], V]{stack: &stack[Comparable[K], V]{p: nil, pp: nil}, node: &valRoot})
	}
}

/*
Gets a 'Just' nodeStack or 'Nothing' if it doesn't exist
*/
func getNodeStackHelp[K, V any](targetKey Comparable[K], ns *nodeStack[Comparable[K], V]) maybe.Maybe[*node[Comparable[K], V]] {
getNodeStackHelpL:
	for {
		switch targetKey.Cmp(ns.node.key) {
		case -1:
			if ns.node.left == nil {
				return maybe.Nothing{}
			} else {
				newStack := &stack[Comparable[K], V]{pp: ns.stack, p: ns.node}
				valL := *ns.node.left
				ns.node.left = &valL
				newNs := &nodeStack[Comparable[K], V]{node: ns.node.left, stack: newStack}
				ns = newNs
				continue getNodeStackHelpL
			}
		case 0:
			return maybe.Just[*nodeStack[Comparable[K], V]]{Value: ns}
		case +1:
			if ns.node.right == nil {
				return maybe.Nothing{}
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
func Get[K, V any](targetKey Comparable[K], d Dict[Comparable[K], V]) maybe.Maybe[V] {
	root := d.rbt().root
	if root == nil {
		return maybe.Nothing{}
	} else {
		return getHelp(targetKey, root)
	}
}

// Determine the number of key-value pairs in the dictionary.
func Size[K, V any](d Dict[Comparable[K], V]) Int {
	root := d.rbt().root

	if root == nil {
		return 0
	}
	queue := []*node[Comparable[K], V]{root}
	count := 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:] // Dequeue
		count++

		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}

	return Int(count)
}

func getHelp[K, V any](targetKey Comparable[K], n *node[Comparable[K], V]) maybe.Maybe[V] {
getHelpL:
	for {
		if n == nil {
			return maybe.Nothing{}
		} else {
			switch targetKey.Cmp(n.key) {
			case -1:
				n = n.left
				continue getHelpL
			case 0:
				return maybe.Just[V]{Value: n.value}
			case +1:
				n = n.right
				continue getHelpL
			}
		}
	}
}

// LISTS

// Get all of the keys in a dictionary, sorted from lowest to highest.
func Keys[K, V any](d Dict[Comparable[K], V]) list.List[Comparable[K]] {
	return Foldr(
		func(key Comparable[K], _ V, keyList list.List[Comparable[K]]) list.List[Comparable[K]] {
			return list.Cons(key, keyList)
		},
		list.Empty[Comparable[K]](),
		d,
	)
}

// Get all of the values in a dictionary, in the order of their keys.
func Values[K, V any](d Dict[Comparable[K], V]) list.List[V] {
	return Foldr(
		func(_ Comparable[K], value V, valueList list.List[V]) list.List[V] {
			return list.Cons(value, valueList)
		},
		list.Empty[V](),
		d,
	)
}

// Convert a dictionary into an association list of key-value pairs, sorted by keys.
func ToList[K, V any](d Dict[Comparable[K], V]) list.List[tuple.Tuple2[Comparable[K], V]] {
	return Foldr(
		func(key Comparable[K], value V, xs list.List[tuple.Tuple2[Comparable[K], V]]) list.List[tuple.Tuple2[Comparable[K], V]] {
			return list.Cons(tuple.Pair(key, value), xs)
		},
		list.Empty[tuple.Tuple2[Comparable[K], V]](),
		d,
	)
}

// Convert an association list into a dictionary.
func FromList[K, V any](l list.List[tuple.Tuple2[Comparable[K], V]]) Dict[Comparable[K], V] {
	return list.Foldl(
		func(t tuple.Tuple2[Comparable[K], V], d Dict[Comparable[K], V]) Dict[Comparable[K], V] {
			return Insert(tuple.First(t), tuple.Second(t), d)
		},
		Empty[K, V](),
		l,
	)
}

// TRANSFORM

// Apply a function to all values in a dictionary.
func Map[K, V, B any](f func(key Comparable[K], value V) B, d Dict[Comparable[K], V]) Dict[Comparable[K], B] {
	if d.rbt().root == nil {
		return Empty[K, B]()
	} else {
		return &dict[Comparable[K], B]{root: mapHelp(f, d.rbt().root)}
	}
}

func mapHelp[K, V, B any](f func(Comparable[K], V) B, t *node[Comparable[K], V]) *node[Comparable[K], B] {
	if t == nil {
		return nil
	} else {
		return &node[Comparable[K], B]{key: t.key, value: f(t.key, t.value), color: t.color, left: mapHelp(f, t.left), right: mapHelp(f, t.right)}
	}
}

// Fold over the key-value pairs in a dictionary from lowest key to highest key.
func Foldl[K, V, B any](f func(Comparable[K], V, B) B, acc B, d Dict[Comparable[K], V]) B {
	return foldlHelp(f, acc, d.rbt().root)
}

func foldlHelp[K, V, B any](f func(Comparable[K], V, B) B, acc B, t *node[Comparable[K], V]) B {
foldlHelpL:
	for {
		if t == nil {
			return acc
		} else {
			var key = t.key
			var value = t.value
			var left = t.left
			var right = t.right
			tempFunc, tempAcc, tempT := f, f(key, value, foldlHelp(f, acc, left)), right
			f = tempFunc
			acc = tempAcc
			t = tempT
			continue foldlHelpL
		}
	}
}

// Fold over the key-value pairs in a dictionary from highest key to lowest key.
func Foldr[K, V, B any](f func(Comparable[K], V, B) B, acc B, d Dict[Comparable[K], V]) B {
	return foldrHelp(f, acc, d.rbt().root)
}

func foldrHelp[K, V, B any](f func(Comparable[K], V, B) B, acc B, t *node[Comparable[K], V]) B {
foldrHelpL:
	for {
		if t == nil {
			return acc
		} else {
			var key = t.key
			var value = t.value
			var left = t.left
			var right = t.right
			tempFunc, tempAcc, tempT := f, f(key, value, foldrHelp(f, acc, right)), left
			f = tempFunc
			acc = tempAcc
			t = tempT
			continue foldrHelpL
		}
	}
}

// Keep only the key-value pairs that pass the given test.
func Filter[K, V any](isGood func(Comparable[K], V) bool, d Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	return Foldl(
		func(k Comparable[K], v V, d Dict[Comparable[K], V]) Dict[Comparable[K], V] {
			if isGood(k, v) {
				return Insert(k, v, d)
			} else {
				return d
			}
		},
		Empty[K, V](),
		d,
	)
}

// Partition a dictionary according to some test. The first dictionary
// contains all key-value pairs which passed the test, and the second contains
// the pairs that did not.
func Partition[K, V any](isGood func(Comparable[K], V) bool, d Dict[Comparable[K], V]) tuple.Tuple2[Dict[Comparable[K], V], Dict[Comparable[K], V]] {
	add := func(
		key Comparable[K],
		value V,
		td tuple.Tuple2[Dict[Comparable[K], V], Dict[Comparable[K], V]],
	) tuple.Tuple2[Dict[Comparable[K], V], Dict[Comparable[K], V]] {
		t1, t2 := tuple.First(td), tuple.Second(td)
		if isGood(key, value) {
			return tuple.Pair(Insert(key, value, t1), t2)
		} else {
			return tuple.Pair(t1, Insert(key, value, t2))
		}
	}
	return Foldl(add, tuple.Pair(Empty[K, V](), Empty[K, V]()), d)
}

// COMBINE

// Combine two dictionaries. If there is a collision, preference is given
// to the first dictionary.
func Union[K, V any](t1 Dict[Comparable[K], V], t2 Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	return Foldl[K, V, Dict[Comparable[K], V]](Insert, t2, t1)
}

// Keep a key-value pair when its key appears in the second dictionary.
// Preference is given to values in the first dictionary.
func Intersect[K, V any](t1 Dict[Comparable[K], V], t2 Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	return Filter(func(k Comparable[K], _ V) bool { return Member(k, t2) }, t1)
}

// Keep a key-value pair when its key does not appear in the second dictionary.
func Diff[K, V any](t1 Dict[Comparable[K], V], t2 Dict[Comparable[K], V]) Dict[Comparable[K], V] {
	return Foldl(
		func(k Comparable[K], _ V, t Dict[Comparable[K], V]) Dict[Comparable[K], V] { return Remove(k, t) },
		t1,
		t2,
	)
}

// The most general way of combining two dictionaries. You provide three
// accumulators for when a given key appears:
// 1. Only in the left dictionary.
// 2. In both dictionaries.
// 3. Only in the right dictionary.
// You then traverse all the keys from lowest to highest, building up whatever
// you want.
func Merge[K, A, B, R any](
	leftStep func(Comparable[K], A, R) R,
	bothStep func(Comparable[K], A, B, R) R,
	rightStep func(Comparable[K], B, R) R,
	leftDict Dict[Comparable[K], A],
	rightDict Dict[Comparable[K], B],
	initialResult R,
) R {
	stepState := func(rKey Comparable[K], rValue B, _v0 tuple.Tuple2[list.List[tuple.Tuple2[Comparable[K], A]], R]) tuple.Tuple2[list.List[tuple.Tuple2[Comparable[K], A]], R] {
	stepStateL:
		for {
			list_ := tuple.First(_v0)
			result := tuple.Second(_v0)
			if list_.Cons() == nil /* empty */ {
				return tuple.Pair(list_, rightStep(rKey, rValue, result))
			} else {
				_v2 := list_.Cons().A
				lKey := tuple.First(_v2)
				lValue := tuple.Second(_v2)
				rest := list_.Cons().B
				if lKey.Cmp(rKey) < 0 {
					tempRKey, tempRValue, temp_v0 := rKey, rValue, tuple.Pair(rest, leftStep(lKey, lValue, result))
					rKey = tempRKey
					rValue = tempRValue
					_v0 = temp_v0
					continue stepStateL
				} else {
					if lKey.Cmp(rKey) > 0 {
						return tuple.Pair(list_, rightStep(rKey, rValue, result))
					} else {
						return tuple.Pair(rest, bothStep(lKey, lValue, rValue, result))
					}
				}
			}
		}
	}
	_v3 := Foldl(stepState, tuple.Pair(ToList(leftDict), initialResult), rightDict)
	leftovers := tuple.First(_v3)
	intermediateResult := tuple.Second(_v3)
	return list.Foldl(
		func(_v4 tuple.Tuple2[Comparable[K], A], result R) R {
			k := tuple.First(_v4)
			v := tuple.Second(_v4)
			return leftStep(k, v, result)
		},
		intermediateResult,
		leftovers,
	)
}
