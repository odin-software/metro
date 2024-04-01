package models

import (
	"errors"
	"slices"
)

type Network[T any] struct {
	vertices     map[string]T
	edges        map[string]map[string][]Vector
	hashFunction func(T) string
}

func NewNetwork[T any](hashF func(T) string) Network[T] {
	return Network[T]{
		vertices:     make(map[string]T),
		edges:        make(map[string]map[string][]Vector),
		hashFunction: hashF,
	}
}

// Inserts a vertex into the graph.
// If the vertex already exists, it returns an error, otherwise it returns nil.
// It also creates an empty map for the edges of the new vertex.
func (gr *Network[T]) InsertVertex(vertex T) error {
	key := gr.hashFunction(vertex)
	if _, ok := gr.vertices[key]; ok {
		return errors.New("this vertex already exists in the graph")
	}

	gr.vertices[key] = vertex
	gr.edges[key] = make(map[string][]Vector)
	return nil
}

// Works like InsertVertex, but for a slice of vertices.
func (gr *Network[T]) InsertVertices2(vertices []T) error {
	for _, v := range vertices {
		err := gr.InsertVertex(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gr *Network[T]) InsertVertices(vertices []*T) error {
	for _, v := range vertices {
		err := gr.InsertVertex(*v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Inserts an edge between two vertices, the points are positions in between the vertices.
func (gr *Network[T]) InsertEdge(firstVertex T, secondVertex T, points []Vector) error {
	firstKey := gr.hashFunction(firstVertex)
	secondKey := gr.hashFunction(secondVertex)

	if _, ok := gr.edges[firstKey]; !ok {
		return errors.New("the first vertex does not exists in the graph")
	}
	if _, ok := gr.edges[secondKey]; !ok {
		return errors.New("the first vertex does not exists in the graph")
	}

	gr.edges[firstKey][secondKey] = points
	rev := make([]Vector, len(points))
	copy(rev, points)
	slices.Reverse(rev)
	gr.edges[secondKey][firstKey] = rev

	return nil
}

func (gr *Network[T]) GetVertexFromKey(key string) (T, error) {
	if _, ok := gr.vertices[key]; !ok {
		return *new(T), errors.New("this vertex doesnt exists in the graph")
	}
	return gr.vertices[key], nil
}

func (gr *Network[T]) GetVertexFromValue(value T) (T, error) {
	key := gr.hashFunction(value)
	vertex, ok := gr.vertices[key]
	if !ok {
		return *new(T), errors.New("this vertex is not on the graph")
	}

	return vertex, nil
}

func (gr *Network[T]) GetVertices() []T {
	values := make([]T, 0, len(gr.vertices))
	for _, v := range gr.vertices {
		values = append(values, v)
	}
	return values
}

func (gr *Network[T]) UpdateVertex(vertex T) error {
	key := gr.hashFunction(vertex)

	if _, ok := gr.vertices[key]; !ok {
		return errors.New("this vertex is not on the graph")
	}
	gr.vertices[key] = vertex

	return nil
}

func (gr *Network[T]) DeleteVertex(vertex T) error {
	key := gr.hashFunction(vertex)
	if _, ok := gr.vertices[key]; !ok {
		return errors.New("this vertex does not exists")
	}

	delete(gr.vertices, key)
	delete(gr.edges, key)

	// Update other edges
	for _, v := range gr.edges {
		delete(v, key)
	}

	return nil
}

func (gr *Network[T]) GetEdges(vertex T) (map[string][]Vector, error) {
	key := gr.hashFunction(vertex)
	v, ok := gr.edges[key]
	if !ok {
		return nil, errors.New("this vertex does not exists in the graph")
	}
	return v, nil
}

func (gr *Network[T]) AreConnected(firstVertex T, secondVertex T) (std []Vector, e error) {
	firstKey := gr.hashFunction(firstVertex)
	secondKey := gr.hashFunction(secondVertex)

	firstMap, ok := gr.edges[firstKey]
	if !ok {
		return std, errors.New("the first vertex does not exists")
	}
	if _, ok := gr.edges[secondKey]; !ok {
		return std, errors.New("the second vertex does not exists")
	}

	connections, ok := firstMap[secondKey]
	if !ok {
		return std, errors.New("these vertices are not connected")
	}

	return connections, nil
}

func (gr *Network[T]) UpdateEdgeValue(firstVertex T, secondVertex T, points []Vector) error {
	firstKey := gr.hashFunction(firstVertex)
	secondKey := gr.hashFunction(secondVertex)

	firstMap, ok := gr.edges[firstKey]
	if !ok {
		return errors.New("the first vertex does not exists")
	}
	secondMap, ok := gr.edges[secondKey]
	if !ok {
		return errors.New("the second vertex does not exists")
	}

	_, ok = firstMap[secondKey]
	if !ok {
		return errors.New("these vertices are not connected")
	}
	_, ok = secondMap[firstKey]
	if !ok {
		return errors.New("these vertices are not connected")
	}

	firstMap[secondKey] = points
	secondMap[firstKey] = points

	return nil
}

func (gr *Network[T]) DeleteEdge(firstVertex T, secondVertex T) error {
	firstKey := gr.hashFunction(firstVertex)
	secondKey := gr.hashFunction(secondVertex)

	firstMap, ok := gr.edges[firstKey]
	if !ok {
		return errors.New("the first vertex does not exists")
	}
	secondMap, ok := gr.edges[secondKey]
	if !ok {
		return errors.New("the second vertex does not exists")
	}

	if _, ok = firstMap[secondKey]; !ok {
		return errors.New("these vertices are not connected")
	}
	delete(firstMap, secondKey)
	delete(secondMap, secondKey)

	return nil
}
