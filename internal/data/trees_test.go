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

func TestFixedInsert(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 1; i < 50; i++ {
		err := tree.Insert(NodeValue[int]{i, 0})
		if err != nil {
			t.Fatal(err)
		}
	}

	if tree.root.right.value.idx == 1 {
		t.Fatal("A balanced tree should not have this number as a root")
	}
}

func TestFixedInsertInverse(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 40; i > 0; i-- {
		err := tree.Insert(NodeValue[int]{i, 0})
		if err != nil {
			t.Fatal(err)
		}
	}

	if tree.root.right.value.idx == 40 {
		t.Fatal("A balanced tree should not have this number as a root")
	}
}

func TestGetFromIndex(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 1; i < 50; i++ {
		err := tree.Insert(NodeValue[int]{i, 0})
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err := tree.Get(448)
	if err == nil {
		t.Fatal("The value 448 should return an error.")
	}

	val, err := tree.Get(22)
	if err != nil {
		t.Fatal("The value 22 should be found in the tree.")
	}

	if val.value.idx != 22 {
		t.Fatal("The wrong value was returned.")
	}
}

func TestGetValueFromIndex(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 1; i < 50; i++ {
		err := tree.Insert(NodeValue[int]{i, 0})
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err := tree.GetValue(448)
	if err == nil {
		t.Fatal("The value 448 should return an error.")
	}

	val, err := tree.GetValue(41)
	if err != nil {
		t.Fatal("The value 22 should be found in the tree.")
	}

	if val.idx != 41 {
		t.Fatal("The wrong value was returned.")
	}
}

func TestCount(t *testing.T) {
	tree := NewTree[int]()

	count := tree.Count()
	if count != 0 {
		t.Fatal("The count should be zero.")
	}

	// Inserting testing data
	for i := 1; i < 50; i++ {
		err := tree.Insert(NodeValue[int]{i, 0})
		if err != nil {
			t.Fatal(err)
		}
	}

	fullCount := tree.Count()
	if fullCount != 49 {
		t.Fatal("The count should be 49.")
	}
}

func TestSmallCount(t *testing.T) {
	tree := NewTree[int]()

	count := tree.Count()
	if count != 0 {
		t.Fatal("The count should be zero.")
	}

	err := tree.Insert(NodeValue[int]{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	partialCount := tree.Count()
	if partialCount != 1 {
		t.Fatal("The count should be 1.")
	}

	// Inserting testing data
	for i := 1; i < 3; i++ {
		err := tree.Insert(NodeValue[int]{i, 0})
		if err != nil {
			t.Fatal(err)
		}
	}

	fullCount := tree.Count()
	if fullCount != 3 {
		t.Fatal("The count should be 3.")
	}
}

func TestGetNodesValue(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 0; i < 30; i++ {
		err := tree.Insert(NodeValue[int]{i, 9})
		if err != nil {
			t.Fatal(err)
		}
	}

	vals := tree.GetNodesValues()
	if len(vals) != 30 {
		t.Fatalf("Expected 30 nodes but got %v", len(vals))
	}
}

func TestGetNodesEmpty(t *testing.T) {
	tree := NewTree[int]()

	vals := tree.GetNodesValues()
	if len(vals) != 0 {
		t.Fatalf("Expected 0 nodes but got %v", len(vals))
	}
}

func TestDeleteNodeWithNoChild(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 0; i < 3; i++ {
		err := tree.Insert(NodeValue[int]{i, 9})
		if err != nil {
			t.Fatal(err)
		}
	}

	if tree.Count() != 3 {
		t.Fatal("The count should be 3.")
	}
	ok := tree.Delete(3)
	if ok {
		t.Fatal("This should have errored because that idx does not exists.")
	}
	ok = tree.Delete(2)
	if !ok {
		t.Fatal("This should not throw an error, since 2 exists and doesn't have children.")
	}
	if tree.Count() != 2 {
		t.Fatal("The deleted node is still being counted.")
	}
}

func TestDeleteNodeWithOneChild(t *testing.T) {
	tree := NewTree[int]()

	// Inserting testing data
	for i := 0; i < 4; i++ {
		err := tree.Insert(NodeValue[int]{i, 9})
		if err != nil {
			t.Fatal(err)
		}
	}

	if tree.Count() != 4 {
		t.Fatal("The count should be 4.")
	}

	// testing one child, right
	ok := tree.Delete(2)
	if !ok {
		t.Fatal("This should not throw an error, since 2 exists and only has one child.")
	}
	if tree.Count() != 3 {
		t.Fatal("The deleted node is still being counted.")
	}

	// testing one child, left
	tree.Insert(NodeValue[int]{-1, 9})
	ok = tree.Delete(0)
	if !ok {
		t.Fatal("This should not throw an error, since 0 exists and only has one child.")
	}
	if tree.Count() != 3 {
		t.Fatal("The deleted node is still being counted.")
	}
}

func TestDeleteNodeWithTwoChildren(t *testing.T) {
	tree := NewTree[int]()
	nodes := 20

	// Inserting testing data
	for i := 0; i < nodes; i++ {
		err := tree.Insert(NodeValue[int]{i, 9})
		if err != nil {
			t.Fatal(err)
		}
	}

	if tree.Count() != nodes {
		t.Fatalf("The count should be %v.", nodes)
	}

	// testing both children
	ok := tree.Delete(17)
	if !ok {
		t.Fatal("This should not throw an error, since 17 exists and has two child.")
	}
	t.Log(tree.Count())
	if tree.Count() != nodes-1 {
		t.Fatal("The deleted node is still being counted.")
	}

	// testing both children
	ok = tree.Delete(3)
	if !ok {
		t.Fatal("This should not throw an error, since 3 exists and has two child.")
	}
	t.Log(tree.Count())
	if tree.Count() != nodes-2 {
		t.Fatal("The deleted node is still being counted.")
	}
}

func TestUpdateValue(t *testing.T) {
	tree := NewTree[int]()

	tree.Insert(NodeValue[int]{
		idx: 2,
		val: 4,
	})
	node, err := tree.Get(2)
	if err != nil {
		t.Fatal("This value is inside the tree.")
	}
	if node.value.val != 4 {
		t.Fatal("The value was saved incorrectly.")
	}

	err = tree.UpdateValue(node.value.idx, 9)
	if err != nil {
		t.Fatal("The value should be in the tree.")
	}
	if node.value.val != 9 {
		t.Fatal("The value was not updated.")
	}
}
