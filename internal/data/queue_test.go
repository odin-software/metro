package data

import "testing"

func TestNewQueue(t *testing.T) {
	q := NewQueue[string]()
	if q.items == nil {
		t.Fatal("The items slice was not created correctly.")
	}
}

func TestQueueValues(t *testing.T) {
	q := NewQueue[string]()
	q.Q("Ciudad 1")
	q.Q("Ciudad 2")
	q.Q("Ciudad 3")

	if q.items[0] != "Ciudad 1" {
		t.Fatal("The item was enqueued incorrectly.")
	}
	if q.items[2] != "Ciudad 3" {
		t.Fatal("The item was enqueued incorrectly.")
	}
}

func TestQueueListOfValues(t *testing.T) {
	q := NewQueue[string]()
	data := []string{
		"Ciudad 1",
		"Ciudad 2",
		"Ciudad 3",
	}
	q.QList(data)

	if q.items[0] != "Ciudad 1" {
		t.Fatal("The item was enqueued incorrectly.")
	}
	if q.items[2] != "Ciudad 3" {
		t.Fatal("The item was enqueued incorrectly.")
	}
}

func TestDequeue(t *testing.T) {
	q := NewQueue[string]()
	q.Q("Ciudad 1")
	q.Q("Ciudad 2")
	q.Q("Ciudad 3")

	if q.items[0] != "Ciudad 1" {
		t.Fatal("The item was enqueued incorrectly.")
	}
	val, err := q.DQ()
	if err != nil {
		t.Fatal("This should have items still.")
	}
	if val != "Ciudad 1" {
		t.Fatal("The item was enqueued incorrectly.")
	}
	_, _ = q.DQ()
	_, _ = q.DQ()
	_, err = q.DQ()
	if err == nil {
		t.Fatal("This should error since there are no more items in the queue.")
	}
}

func TestQueueSize(t *testing.T) {
	q := NewQueue[string]()
	q.Q("Ciudad 1")
	q.Q("Ciudad 2")
	q.Q("Ciudad 3")

	if q.Size() != 3 {
		t.Fatal("The size is incorrect.")
	}
}
