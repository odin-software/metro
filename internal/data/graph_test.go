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
	g := NewGraph[TestStruct](tsHashFuncion)

	if g.vertices == nil {
		t.Fatal("The graph vertices have not been created.")
	}
	if g.edges == nil {
		t.Fatal("The graph edges have not been created.")
	}
}

func TestInsertVertex(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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
	if g.vertices[ts1.name].name != ts1.name {
		t.Fatal("The vertex was added with wrong information.")
	}
	if g.edges[ts2.name] == nil {
		t.Fatal("Failed on edgelist creation.")
	}
	if g.edges[ts2.name] == nil {
		t.Fatal("Failed on edgelist creation.")
	}

	err := g.InsertVertex(ts1)
	if err == nil {
		t.Fatal("It should error on duplicated vertex insertion.")
	}
}

func TestInsertEdge(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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

	err := g.InsertEdge(ts1, ts2, 0)
	if err == nil {
		t.Fatal("You cannot insert a non-positive edge weight.")
	}
	err = g.InsertEdge(ts1, ts2, 4)
	if err != nil {
		t.Fatal(err)
	}

	if g.edges[g.hashFunction(ts1)][g.hashFunction(ts2)] != 4 {
		t.Fatal("Edges have not been created.")
	}
}

func TestNonExistingVertexInsertEdge(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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

func TestGetVertexFromKey(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}

	g.InsertVertex(ts1)

	vertex, err := g.GetVertexFromKey(ts1.name)
	if err != nil {
		t.Fatal(err)
	}
	if vertex.color != "#223332" {
		t.Fatal("Got the wrong value.")
	}

	vertex, err = g.GetVertexFromKey(ts2.name)
	if err == nil {
		t.Fatal("This key should not be in the graph.")
	}
}
func TestGetVertexFromValue(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}

	g.InsertVertex(ts1)

	vertex, err := g.GetVertexFromValue(ts1)
	if err != nil {
		t.Fatal(err)
	}
	if vertex.color != "#223332" {
		t.Fatal("Got the wrong value.")
	}

	vertex, err = g.GetVertexFromValue(ts2)
	if err == nil {
		t.Fatal("This key should not be in the graph.")
	}
}

func TestGetVertices(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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
}

func TestUpdateVertex(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

	ts1 := TestStruct{
		name:  "Ciudad 1",
		color: "#223332",
	}
	ts2 := TestStruct{
		name:  "Ciudad 2",
		color: "#225332",
	}

	g.InsertVertex(ts1)

	ts1.color = "#112112"
	err := g.UpdateVertex(ts1)
	if err != nil {
		t.Fatal(err)
	}
	ts1Updated, err := g.GetVertexFromKey(ts1.name)
	if err != nil {
		t.Fatal(err)
	}
	if ts1Updated.color != "#112112" {
		t.Fatal("The color was not updated.")
	}

	err = g.UpdateVertex(ts2)
	if err == nil {
		t.Fatal("It should fail because the vertex does not exist.")
	}
}

func TestDeleteVertex(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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

	err := g.DeleteVertex(ts2)
	if err != nil {
		t.Fatal(err)
	}

	vertices = g.GetVertices()
	if len(vertices) != 2 {
		t.Fatal("The vertices list was not updated.")
	}

	err = g.DeleteVertex(ts2)
	if err == nil {
		t.Fatal("This should fail because the vertex was deleted.")
	}
}

func TestGetEdges(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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
	ts4 := TestStruct{
		name:  "Ciudad Inexistente",
		color: "#2e3a22",
	}

	g.InsertVertex(ts1)
	g.InsertVertex(ts2)
	g.InsertVertex(ts3)

	emptyEdges, err := g.GetEdges(ts2)
	if err != nil {
		t.Fatal(err)
	}
	if len(emptyEdges) != 0 {
		t.Fatal("There should be no edges.")
	}

	err = g.InsertEdge(ts1, ts2, 4)
	if err != nil {
		t.Fatal(err)
	}
	err = g.InsertEdge(ts1, ts3, 7)
	if err != nil {
		t.Fatal(err)
	}
	err = g.InsertEdge(ts2, ts3, 9)
	if err != nil {
		t.Fatal(err)
	}

	edges, err := g.GetEdges(ts1)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 2 {
		t.Fatal("Ts1 should have two edges connected.")
	}
	if edges[g.hashFunction(ts2)] != 4 {
		t.Fatal("The connection value is incorrect.")
	}
	edges, err = g.GetEdges(ts2)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 2 {
		t.Fatal("Ts2 should have two edges connected.")
	}
	if edges[g.hashFunction(ts3)] != 9 {
		t.Fatal("The connection value is incorrect.")
	}

	_, err = g.GetEdges(ts4)
	if err == nil {
		t.Fatal("This should error bacause the vertex does not exist.")
	}
}

func TestUpdateEdgeValue(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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
	if values[g.hashFunction(ts1)] != 8 {
		t.Fatal("the value is not correct")
	}

	err = g.UpdateEdgeValue(ts1, ts2, 3)
	if err != nil {
		t.Fatal(err)
	}
	values, err = g.GetEdges(ts2)
	if err != nil {
		t.Fatal(err)
	}
	if values[g.hashFunction(ts1)] != 3 {
		t.Fatal("the value is not correct")
	}

	err = g.UpdateEdgeValue(ts1, ts3, 4)
	if err == nil {
		t.Fatal("This should fail because these vertices are not connected.")
	}
}

func TestAreConnected(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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

	weight, err := g.AreConnected(ts1, ts2)
	if err != nil {
		t.Fatal("These two vertices should be connected.")
	}
	if weight != 8 {
		t.Fatal("Wrong value in the connection.")
	}
	weight, err = g.AreConnected(ts2, ts3)
	if err != nil {
		t.Fatal("These two vertices are connected.")
	}
	if weight != 3 {
		t.Fatal("Wrong value in the connection.")
	}

	// opposite way should work too
	weight, err = g.AreConnected(ts3, ts2)
	if err != nil {
		t.Fatal("These two vertices are connected.")
	}
	if weight != 3 {
		t.Fatal("Wrong value in the connection.")
	}

	// testing unconnected
	_, err = g.AreConnected(ts1, ts3)
	if err == nil {
		t.Fatal("These two vertices are not connected.")
	}
}

func TestDeleteEdge(t *testing.T) {
	g := NewGraph[TestStruct](tsHashFuncion)

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

	_, err := g.AreConnected(ts1, ts2)
	if err != nil {
		t.Fatal("Expected this to to be connected.")
	}

	err = g.DeleteEdge(ts1, ts2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.AreConnected(ts1, ts2)
	if err == nil {
		t.Fatal("The vertices should not be connected.")
	}
}
