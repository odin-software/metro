package data

import (
	"errors"
)

type Graph[T any, V any] struct {
	vertices     []T
	edges        []RBTree[V]
	hash         map[string]int
	hashFunction func(T) string
}

func NewGraph[T any, V any](hashF func(T) string) Graph[T, V] {
	return Graph[T, V]{
		vertices:     []T{},
		edges:        []RBTree[V]{},
		hash:         map[string]int{},
		hashFunction: hashF,
	}
}

func (gr *Graph[T, V]) InsertVertex(val T) {
	gr.vertices = append(gr.vertices, val)

	gr.edges = append(gr.edges, RBTree[V]{})
	gr.hash[gr.hashFunction(val)] = len(gr.vertices) - 1
}

func (gr *Graph[T, V]) InsertEdge(firstVertex T, secondVertex T, weight V) error {
	fIdx, ok := gr.hash[gr.hashFunction(firstVertex)]
	if !ok {
		return errors.New("The first vertex does not exists on the graph.")
	}
	sIdx, ok := gr.hash[gr.hashFunction(secondVertex)]
	if !ok {
		return errors.New("The second vertex does not exists on the graph.")
	}

	err := gr.edges[fIdx].Insert(NodeValue[V]{
		idx: sIdx,
		val: weight,
	})
	if err != nil {
		return errors.New("This edge already exists")
	}
	err2 := gr.edges[sIdx].Insert(NodeValue[V]{
		idx: fIdx,
		val: weight,
	})
	if err2 != nil {
		return errors.New("This edge already exists")
	}

	return nil
}

func (gr *Graph[T, V]) GetEdges(v T) ([]NodeValue[V], error) {
	idx, ok := gr.hash[gr.hashFunction(v)]
	if !ok {
		return nil, errors.New("This vertex does not exists on the graph.")
	}

	edges := gr.edges[idx].GetNodesValues()
	return edges, nil
}
