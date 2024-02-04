package data

import (
	"errors"
	"reflect"
)

type NodeValue[T any] struct {
	idx int
	val T
}

type RBNode[V any] struct {
	red    bool
	parent *RBNode[V]
	value  NodeValue[V]
	left   *RBNode[V]
	right  *RBNode[V]
}

type RBTree[V any] struct {
	leaf *RBNode[V]
	root *RBNode[V]
}

func NewTree[V any]() *RBTree[V] {
	leaf := RBNode[V]{}
	return &RBTree[V]{
		leaf: &leaf,
		root: &leaf,
	}
}

func (rbt *RBTree[T]) Insert(val NodeValue[T]) error {
	newNode := RBNode[T]{
		red:   true,
		value: val,
		left:  rbt.leaf,
		right: rbt.leaf,
	}

	var parent *RBNode[T]
	current := rbt.root
	for current != rbt.leaf {
		parent = current
		if newNode.value.idx < current.value.idx {
			current = current.left
		} else if newNode.value.idx > current.value.idx {
			current = current.right
		} else {
			return errors.New("RBTree: Inserting a duplicated index")
		}
	}

	newNode.parent = parent
	if parent == nil {
		rbt.root = &newNode
	} else if newNode.value.idx < parent.value.idx {
		parent.left = &newNode
	} else {
		parent.right = &newNode
	}

	rbt.fix_insert(&newNode)

	return nil
}

func (rbt *RBTree[T]) rotate_left(x *RBNode[T]) {
	if x == rbt.leaf || x.right == rbt.leaf {
		return
	}

	y := x.right
	x.right = y.left
	if y.left != rbt.leaf {
		y.left.parent = x
	}

	y.parent = x.parent
	if x.parent == nil {
		rbt.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func (rbt *RBTree[T]) rotate_right(x *RBNode[T]) {
	if x == rbt.leaf || x.left == rbt.leaf {
		return
	}

	y := x.left
	x.left = y.right
	if y.right != rbt.leaf {
		y.right.parent = x
	}

	y.parent = x.parent
	if x.parent == nil {
		rbt.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y
}

func (rbt *RBTree[V]) fix_insert(newNode *RBNode[V]) {
	for newNode != rbt.root && newNode.parent.red {
		if newNode.parent == newNode.parent.parent.right {
			uncle := newNode.parent.parent.left
			if uncle.red {
				uncle.red = false
				newNode.parent.red = false
				newNode.parent.parent.red = true
				newNode = newNode.parent.parent
			} else {
				if newNode == newNode.parent.left {
					newNode = newNode.parent
					rbt.rotate_right(newNode)
				}
				newNode.parent.red = false
				newNode.parent.parent.red = true
				rbt.rotate_left(newNode.parent.parent)
			}
		} else {
			uncle := newNode.parent.parent.right
			if uncle.red {
				uncle.red = false
				newNode.parent.red = false
				newNode.parent.parent.red = true
				newNode = newNode.parent.parent
			} else {
				if newNode == newNode.parent.right {
					newNode = newNode.parent
					rbt.rotate_left(newNode)
				}
				newNode.parent.red = false
				newNode.parent.parent.red = true
				rbt.rotate_right(newNode.parent.parent)
			}
		}
	}
	rbt.root.red = false
}

func (rbt *RBTree[V]) Get(idx int) (*RBNode[V], error) {
	current := rbt.root
	for current != rbt.leaf && idx != current.value.idx {
		if idx < current.value.idx {
			current = current.left
		} else {
			current = current.right
		}
	}

	if reflect.DeepEqual(current, &RBNode[V]{}) {
		return nil, errors.New("This value is not in this tree.")
	}
	return current, nil
}

func (rbt *RBTree[V]) Count() int {
	return countNodes[V](rbt.root)
}

func countNodes[V any](root *RBNode[V]) int {
	if reflect.DeepEqual(root, &RBNode[V]{}) {
		return 0
	}

	return 1 + countNodes(root.left) + countNodes(root.right)
}
