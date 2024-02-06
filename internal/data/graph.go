package data

import (
	"errors"
)

type Graph[T comparable, V any] struct {
	vertices     []T
	edges        []RBTree[V]
	hash         map[string]int
	hashFunction func(T) string
}

func NewGraph[T comparable, V any](hashF func(T) string) Graph[T, V] {
	return Graph[T, V]{
		vertices:     []T{},
		edges:        []RBTree[V]{},
		hash:         map[string]int{},
		hashFunction: hashF,
	}
}

func (gr *Graph[T, V]) InsertVertex(val T) error {
	for _, v := range gr.vertices {
		if v == val {
			return errors.New("duplicated vertex")
		}
	}

	gr.vertices = append(gr.vertices, val)
	gr.edges = append(gr.edges, *NewTree[V]())
	gr.hash[gr.hashFunction(val)] = len(gr.vertices) - 1
	return nil
}

func (gr *Graph[T, V]) InsertEdge(firstVertex T, secondVertex T, weight V) error {
	fIdx, ok := gr.hash[gr.hashFunction(firstVertex)]
	if !ok {
		return errors.New("the first vertex does not exists on the graph")
	}
	sIdx, ok := gr.hash[gr.hashFunction(secondVertex)]
	if !ok {
		return errors.New("the second vertex does not exists on the graph")
	}

	err := gr.edges[fIdx].Insert(NodeValue[V]{
		idx: sIdx,
		val: weight,
	})
	if err != nil {
		return errors.New("this edge already exists")
	}
	err2 := gr.edges[sIdx].Insert(NodeValue[V]{
		idx: fIdx,
		val: weight,
	})
	if err2 != nil {
		return errors.New("this edge already exists")
	}

	return nil
}

func (gr *Graph[T, V]) GetEdges(v T) ([]NodeValue[V], error) {
	idx, ok := gr.hash[gr.hashFunction(v)]
	if !ok {
		return nil, errors.New("this vertex does not exists on the graph")
	}

	edges := gr.edges[idx].GetNodesValues()
	return edges, nil
}

func (gr *Graph[T, V]) AreConnected(firstVertex T, secondVertex T) (NodeValue[V], bool) {
	idx, ok := gr.hash[gr.hashFunction(secondVertex)]
	if !ok {
		return NodeValue[V]{}, false
	}
	values, err := gr.GetEdges(firstVertex)
	if err != nil {
		return NodeValue[V]{}, false
	}

	for _, v := range values {
		if v.idx == idx {
			return v, true
		}
	}

	return NodeValue[V]{}, false
}

func (gr *Graph[T, V]) DeleteEdge(firstVertex T, secondVertex T) error {
	fIdx, ok := gr.hash[gr.hashFunction(firstVertex)]
	if !ok {
		return errors.New("the first vertex does not exists on the graph")
	}
	sIdx, ok := gr.hash[gr.hashFunction(secondVertex)]
	if !ok {
		return errors.New("the second vertex does not exists on the graph")
	}
	ok = gr.edges[fIdx].Delete(sIdx)
	if !ok {
		return errors.New("this connection does not exists")
	}
	ok = gr.edges[sIdx].Delete(fIdx)
	if !ok {
		return errors.New("this connection does not exists")
	}

	return nil
}

func (gr *Graph[T, V]) GetVertices() []T {
	return gr.vertices
}

func (gr *Graph[T, V]) GetVertex(v T) (*T, error) {
	idx, ok := gr.hash[gr.hashFunction(v)]
	if !ok {
		return nil, errors.New("this vertex does not exists on the graph")
	}
	return &gr.vertices[idx], nil
}

func (gr *Graph[T, V]) UpdateEdgeValue(firstVertex T, secondVertex T, weight V) error {
	fIdx, ok := gr.hash[gr.hashFunction(firstVertex)]
	if !ok {
		return errors.New("the first vertex does not exists on the graph")
	}
	sIdx, ok := gr.hash[gr.hashFunction(secondVertex)]
	if !ok {
		return errors.New("the second vertex does not exists on the graph")
	}

	err := gr.edges[fIdx].UpdateValue(sIdx, weight)
	if err != nil {
		return errors.New("this connection does not exists")
	}
	err = gr.edges[sIdx].UpdateValue(fIdx, weight)
	if err != nil {
		return errors.New("this connection does not exists")
	}

	return nil
}
