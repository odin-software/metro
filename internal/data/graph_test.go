package data

import (
	"testing"
)

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
	if g.edges[1].root == nil {
		t.Fatal("Failed on edgelist creation.")
	}
	if g.hash[ts2.name] != 1 {
		t.Fatal("Hash is returning wrong index.")
	}

	err := g.InsertVertex(ts1)
	if err == nil {
		t.Fatal("It should error on duplicated vertex insertion.")
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
		t.Fatal("This should error because the edge should exist already.")
	}
}

func TestNonExistingVertexInsertEdge(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)

	err := g.InsertEdge(ts1, ts2, 4)
	if err != nil {
		t.Fatal(err)
	}

	nonExistingErr := g.InsertEdge(ts2, ts3, 8)
	if nonExistingErr == nil {
		t.Fatal("This should error because the vertex does not exists.")
	}
}

func TestGetEdges(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	emptyEdges, errEmpty := g.GetEdges(ts2)
	if errEmpty != nil {
		t.Fatal(errEmpty)
	}
	if len(emptyEdges) != 0 {
		t.Fatal("There should be no edges.")
	}

	err := g.InsertEdge(ts1, ts2, 4)
	if err != nil {
		t.Fatal(err)
	}
	err2 := g.InsertEdge(ts1, ts3, 7)
	if err2 != nil {
		t.Fatal(err2)
	}
	err3 := g.InsertEdge(ts2, ts3, 9)
	if err3 != nil {
		t.Fatal(err3)
	}

	edges, err := g.GetEdges(ts1)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 2 {
		t.Fatal("Ts1 should have two edges connected.")
	}
	edgests2, errts2 := g.GetEdges(ts2)
	if errts2 != nil {
		t.Fatal(errts2)
	}
	if len(edgests2) != 2 {
		t.Fatal("Ts2 should have two edges connected.")
	}
}

func TestGetVertices(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	vertices := g.GetVertices()
	if len(vertices) != 3 {
		t.Fatal("The vertices list was not returned.")
	}
	if vertices[0].name != ts1.name {
		t.Fatal("The vertices list was not returned.")
	}
	if vertices[1].name != ts2.name {
		t.Fatal("The vertices list was not returned.")
	}
	if vertices[2].name != ts3.name {
		t.Fatal("The vertices list was not returned.")
	}
}

func TestGetVertex(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	v, err := g.GetVertex(ts2)
	if err != nil {
		t.Fatal(err)
	}
	if v.name != ts2.name {
		t.Fatal("The vertex was not returned.")
	}

	_, err = g.GetVertex(TestStruct{
		name:  "johnson",
		color: "azul",
	})
	if err == nil {
		t.Fatal("Searching for a vertex that does not exists should give an error.")
	}
}

func TestUpdateEdgeValue(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	err := g.InsertEdge(ts1, ts2, 8)
	if err != nil {
		t.Fatal(err)
	}
	err = g.InsertEdge(ts2, ts3, 8)
	if err != nil {
		t.Fatal(err)
	}

	values, err := g.GetEdges(ts2)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].val != 8 {
		t.Fatal("the value is not correct")
	}

	g.UpdateEdgeValue(ts1, ts2, 3)
	values, err = g.GetEdges(ts2)
	if err != nil {
		t.Fatal(err)
	}
	if values[0].val != 3 {
		t.Fatal("the value is not correct")
	}
}

func TestAreConnected(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	err := g.InsertEdge(ts1, ts2, 8)
	if err != nil {
		t.Fatal(err)
	}
	err = g.InsertEdge(ts2, ts3, 3)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := g.AreConnected(ts1, ts2)
	if !ok {
		t.Fatal("These two vertices are connected.")
	}
	if val.val != 8 {
		t.Fatal("Wrong value in the connection.")
	}
	val, ok = g.AreConnected(ts2, ts3)
	if !ok {
		t.Fatal("These two vertices are connected.")
	}
	if val.val != 3 {
		t.Fatal("Wrong value in the connection.")
	}
	// opposite way should work too
	val, ok = g.AreConnected(ts3, ts2)
	if !ok {
		t.Fatal("These two vertices are connected.")
	}
	if val.val != 3 {
		t.Fatal("Wrong value in the connection.")
	}
	// testing unconnected
	val, ok = g.AreConnected(ts1, ts3)
	if ok {
		t.Fatal("These two vertices are not connected.")
	}
}

func TestDeleteEdge(t *testing.T) {
	g := NewGraph[TestStruct, int](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}
	ts3 := TestStruct{
		name:  "Ciudad 3",
		color: "#299123",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	g.InsertEdge(ts1, ts3, 4)
	g.InsertEdge(ts1, ts2, 4)

	_, ok := g.AreConnected(ts1, ts2)
	if !ok {
		t.Fatal("Expected this to to be connected.")
	}

	err := g.DeleteEdge(ts1, ts2)
	if err != nil {
		t.Fatal(err)
	}

	_, ok = g.AreConnected(ts1, ts2)
	if ok {
		t.Fatal("The vertices should not be connected.")
	}
}
