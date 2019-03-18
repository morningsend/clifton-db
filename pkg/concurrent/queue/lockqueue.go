package queue

import "sync"

type node struct {
	value interface{}
	next  *node
}

type LockQueue struct {
	head *node
	tail *node
	size int

	headLock sync.Mutex
	tailLock sync.Mutex
}

func NewLockQueue() Queue {
	sentinel := &node{
		value: nil,
		next:  nil,
	}

	return &LockQueue{
		head: sentinel,
		tail: sentinel,
	}
}

func (q *LockQueue) Len() int {
	return q.size
}

func (q *LockQueue) Enqueue(value interface{}) bool {
	newNode := &node{
		next:  nil,
		value: value,
	}
	q.tailLock.Lock()
	defer q.tailLock.Unlock()

	q.tail.next = newNode
	q.tail = newNode
	q.size++

	return false
}

func (q *LockQueue) Dequeue() (value interface{}, ok bool) {
	q.headLock.Lock()
	defer q.headLock.Unlock()

	node := q.head
	newHead := node.next

	if newHead == nil {
		return nil, false
	}

	value = newHead.value
	q.head = newHead

	return value, true
}
