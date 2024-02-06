package data

import (
	"errors"
)

type Graph[T comparable] struct {
	vertices     map[string]T
	edges        map[string]map[string]float64
	hashFunction func(T) string
}

func NewGraph[T comparable](hashF func(T) string) Graph[T] {
	return Graph[T]{
		vertices:     make(map[string]T),
		edges:        make(map[string]map[string]float64),
		hashFunction: hashF,
	}
}

func (gr *Graph[T]) InsertVertex(vertex T) error {
	key := gr.hashFunction(vertex)
	if _, ok := gr.vertices[key]; ok {
		return errors.New("this vertex already exists in the graph")
	}

	gr.vertices[key] = vertex
	gr.edges[key] = make(map[string]float64)
	return nil
}

func (gr *Graph[T]) InsertEdge(firstVertex T, secondVertex T, weight float64) error {
	firstKey := gr.hashFunction(firstVertex)
	secondKey := gr.hashFunction(secondVertex)

	if weight <= 0 {
		return errors.New("the weight should be greater than zero")
	}

	if _, ok := gr.edges[firstKey]; !ok {
		return errors.New("the first vertex does not exists in the graph")
	}
	if _, ok := gr.edges[secondKey]; !ok {
		return errors.New("the first vertex does not exists in the graph")
	}
	gr.edges[firstKey][secondKey] = weight
	gr.edges[secondKey][firstKey] = weight

	return nil
}

func (gr *Graph[T]) GetVertexFromKey(key string) (T, error) {
	if _, ok := gr.vertices[key]; !ok {
		return *new(T), errors.New("this vertex doesnt exists in the graph")
	}
	return gr.vertices[key], nil
}

func (gr *Graph[T]) GetVertexFromValue(value T) (T, error) {
	key := gr.hashFunction(value)
	vertex, ok := gr.vertices[key]
	if !ok {
		return *new(T), errors.New("this vertex is not on the graph")
	}

	return vertex, nil
}

func (gr *Graph[T]) GetVertices() []T {
	values := make([]T, 0, len(gr.vertices))
	for _, v := range gr.vertices {
		values = append(values, v)
	}
	return values
}

func (gr *Graph[T]) UpdateVertex(vertex T) error {
	key := gr.hashFunction(vertex)

	if _, ok := gr.vertices[key]; !ok {
		return errors.New("this vertex is not on the graph")
	}
	gr.vertices[key] = vertex

	return nil
}

func (gr *Graph[T]) DeleteVertex(vertex T) error {
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

func (gr *Graph[T]) GetEdges(vertex T) (map[string]float64, error) {
	key := gr.hashFunction(vertex)
	v, ok := gr.edges[key]
	if !ok {
		return nil, errors.New("this vertex does not exists in the graph")
	}
	return v, nil
}

func (gr *Graph[T]) AreConnected(firstVertex T, secondVertex T) (float64, error) {
	firstKey := gr.hashFunction(firstVertex)
	secondKey := gr.hashFunction(secondVertex)

	firstMap, ok := gr.edges[firstKey]
	if !ok {
		return 0, errors.New("the first vertex does not exists")
	}
	if _, ok := gr.edges[secondKey]; !ok {
		return 0, errors.New("the second vertex does not exists")
	}

	weight, ok := firstMap[secondKey]
	if !ok {
		return 0, errors.New("these vertices are not connected")
	}

	return weight, nil
}

func (gr *Graph[T]) UpdateEdgeValue(firstVertex T, secondVertex T, weight float64) error {
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

	firstMap[secondKey] = weight
	secondMap[firstKey] = weight

	return nil
}

func (gr *Graph[T]) DeleteEdge(firstVertex T, secondVertex T) error {
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
