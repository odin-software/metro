package models

import "errors"

type Queue[T comparable] struct {
	items []T
}

func NewQueue[T comparable]() Queue[T] {
	return Queue[T]{
		items: []T{},
	}
}

func (q *Queue[T]) Q(value T) {
	q.items = append(q.items, value)
}

func (q *Queue[T]) QList(value []T) {
	q.items = append(q.items, value...)
}

func (q *Queue[T]) DQ() (T, error) {
	if len(q.items) == 0 {
		return *new(T), errors.New("there are no more items in the queue")
	}
	val := q.items[0]
	q.items = q.items[1:]
	return val, nil
}

// Deletes all the items from the queue.
func (q *Queue[T]) Clear() {
	if len(q.items) == 0 {
		return
	}

	q.items = []T{}
}

func (q *Queue[T]) Peek() (T, error) {
	if len(q.items) == 0 {
		return *new(T), errors.New("there are no items in the queue")
	}
	return q.items[0], nil
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}
