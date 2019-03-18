package queue

import "testing"

func TestConcurrentLinkedListQueue_Enqueue(t *testing.T) {
	q := NewConcurrentLinkedListQueue()
	const (
		a int = iota
		b
		c
	)
	q.Enqueue(a)
	q.Enqueue(b)
	q.Enqueue(c)
	if q.Len() != 3 {
		t.Error("queue len should 3")
	}
}

func TestConcurrentLinkedListQueue_Dequeue_FIFO(t *testing.T) {
	q := NewConcurrentLinkedListQueue()
	const (
		a int = iota
		b
		c
	)

	q.Enqueue(a)
	q.Enqueue(b)
	q.Enqueue(c)

	x, ok := q.Dequeue()
	if !ok {
		t.Error("should be able to dequeue")
	}
	if x != a {
		t.Errorf("element order not fifo, expect %v got %v", a, x)
	}

	x, ok = q.Dequeue()
	if x != b {
		t.Errorf("element order not fifo, expect %v got %v", b, x)
	}

	x, ok = q.Dequeue()
	if x != c {
		t.Errorf("element order not fifo, expect %v got %v", c, x)
	}

	x, ok = q.Dequeue()
	if ok {
		t.Errorf("should be empty but dequeued element")
	}

	if q.Len() != 0 {
		t.Error("empty queue should be 0")
	}
}
