package queue

import "testing"

func TestLockQueue_Enqueue(t *testing.T) {
	const (
		a = iota
		b
		c
	)

	q := NewLockQueue()

	q.Enqueue(a)
	q.Enqueue(b)
	q.Enqueue(c)

	if q.Len() != 3 {
		t.Errorf("expect len to be %d", 3)
		return
	}
}

func TestLockQueue_Dequeue(t *testing.T) {
	const (
		a = iota
		b
		c
	)

	q := NewLockQueue()
	q.Enqueue(a)
	q.Enqueue(b)
	q.Enqueue(c)

	x, ok := q.Dequeue()

	if !ok {
		t.Errorf("queue should have element to dequeue")
	}

	if x != a {
		t.Errorf("expect FIFO %v but got %v", a, x)
	}

	x, ok = q.Dequeue()
	if x != b {
		t.Errorf("expect FIFO %v got got %v", b, x)
	}

	x, ok = q.Dequeue()
	if x != c {
		t.Errorf("expect FIFO %v got got %v", c, x)
	}

	x, ok = q.Dequeue()
	if ok {
		t.Errorf("empty dequeuing empty queue should NOT be ok")
	}
}
