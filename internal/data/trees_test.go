package data

import "testing"

func TestTreeCreation(t *testing.T) {
	tree := NewTree[string]()
	if tree.leaf.value.val != "" {
		t.Fatal("The tree was created incorrectly")
	}
	if tree.root.value.val != "" {
		t.Fatal("The tree was created incorrectly")
	}
}

func TestInsertNormal(t *testing.T) {
	tree := NewTree[int]()

	edges := [3]NodeValue[int]{
		{
			idx: 4,
			val: 3,
		},
		{
			idx: 9,
			val: 2,
		},
		{
			idx: 2,
			val: 3,
		},
	}

	for _, v := range edges {
		err := tree.Insert(v)
		if err != nil {
			t.Fatal(err)
		}
	}

	if tree.root.right.value.idx != 9 {
		t.Fatal("The inserted value 9 is not on the tree.")
	}
	if tree.root.left.value.idx != 2 {
		t.Fatal("The inserted value 2 is not on the tree.")
	}
}

func TestDuplicateInsertion(t *testing.T) {
	tree := NewTree[int]()

	edges := [3]NodeValue[int]{
		{
			idx: 4,
			val: 3,
		},
		{
			idx: 4,
			val: 2,
		},
		{
			idx: 2,
			val: 3,
		},
	}

	for idx, v := range edges {
		err := tree.Insert(v)
		if idx == 1 && err == nil {
			t.Fatal(err)
		}
	}

	if tree.root.left.value.idx != 2 {
		t.Fatal("The inserted value 2 is not on the tree.")
	}
}

func TestRotateLeft(t *testing.T) {
	tree := NewTree[int]()

	edges := [6]NodeValue[int]{
		{
			idx: 4,
			val: 3,
		},
		{
			idx: 9,
			val: 2,
		},
		{
			idx: 2,
			val: 3,
		},
		{
			idx: 11,
			val: 3,
		},
		{
			idx: 7,
			val: 3,
		},
		{
			idx: 18,
			val: 3,
		},
	}

	for _, v := range edges {
		err := tree.Insert(v)
		if err != nil {
			t.Fatal(err)
		}
	}

	pivot := tree.root.right
	if pivot.value.idx != 9 {
		t.Fatal("The inserted value 9 is not on the tree.")
	}

	// Calling left rotation.
	tree.rotate_left(pivot)

	if tree.root.right.value.idx != 11 {
		t.Fatal("The inserted value 11 is not on the right place.")
	}
	if tree.root.right.right.value.idx != 18 {
		t.Fatal("The inserted value 18 is not on the right place.")
	}
}

func TestRotateRight(t *testing.T) {
	tree := NewTree[int]()

	edges := [6]NodeValue[int]{
		{
			idx: 4,
			val: 3,
		},
		{
			idx: 9,
			val: 2,
		},
		{
			idx: 2,
			val: 3,
		},
		{
			idx: 11,
			val: 3,
		},
		{
			idx: 7,
			val: 3,
		},
		{
			idx: 18,
			val: 3,
		},
	}

	for _, v := range edges {
		err := tree.Insert(v)
		if err != nil {
			t.Fatal(err)
		}
	}

	pivot := tree.root.right
	if pivot.value.idx != 9 {
		t.Fatal("The inserted value 9 is not on the tree.")
	}

	// Calling left rotation.
	tree.rotate_right(pivot)

	if tree.root.right.value.idx != 7 {
		t.Fatal("The inserted value 7 is not on the right place.")
	}
	if tree.root.right.right.value.idx != 9 {
		t.Fatal("The inserted value 9 is not on the right place.")
	}
}
