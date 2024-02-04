package data

import "testing"

type TestStruct struct {
	name  string
	color string
}

var tsHashFuncion = func(ts TestStruct) string {
	return ts.name
}

func TestCreatingAGraph(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	if g.vertices == nil {
		t.Fatal("The graph vertices have not been created.")
	}
	if g.edges == nil {
		t.Fatal("The graph edges have not been created.")
	}
}

func TestInsertVertex(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)

	if len(g.vertices) != 2 {
		t.Fatal("The vertices list was not updated.")
	}
	if len(g.edges) != 2 {
		t.Fatal("The edges list was not updated.")
	}
	if g.vertices[0].name != ts1.name {
		t.Fatal("The vertex was added with wrong information.")
	}
	if g.edges[1].root != nil {
		t.Fatal("Failed on edgelist creation.")
	}
	if g.hash[ts2.name] != 1 {
		t.Fatal("Hash is returning wrong index.")
	}
}

func TestInsertEdge(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)

	err := g.InsertEdge(ts1, ts2, 4)
	if err != nil {
		t.Fatal(err)
	}

	if g.edges[0].Count() != 1 && g.edges[1].Count() != 1 {
		t.Fatal("Edges have not been created.")
	}

	errEdge := g.InsertEdge(ts1, ts2, 3)
	if errEdge == nil {
		t.Fatal("This should error!")
	}
}
