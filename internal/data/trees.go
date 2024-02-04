package data

import "errors"

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

	return nil
}
