package queue

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)


type linkedListNode struct {
	value   interface{}
	nextPtr *linkedListNodePointer
}

type linkedListNodePointer struct {
	node  *linkedListNode
	count int
}

type ConcurrentLinkedListQueue struct {
	head *linkedListNodePointer
	tail *linkedListNodePointer

	sentinel *linkedListNode
}

func NewConcurrentLinkedListQueue() Queue {
	sentinel := &linkedListNode{
		value: nil,
		nextPtr: &linkedListNodePointer{
			node:  nil,
			count: 0,
		},
	}

	return &ConcurrentLinkedListQueue{
		head: &linkedListNodePointer{
			node:  sentinel,
			count: 0,
		},
		tail: &linkedListNodePointer{
			node:  sentinel,
			count: 0,
		},
		sentinel: sentinel,
	}
}

func (q *ConcurrentLinkedListQueue) Enqueue(value interface{}) (ok bool) {
	newNode := &linkedListNode{
		value: value,
		nextPtr: &linkedListNodePointer{
			node:  nil,
			count: 0,
		},
	}
	var tail *linkedListNodePointer
	var nextPtr *linkedListNodePointer

	fmt.Println("inserting", value)
	for {
		tail = q.tail
		nextPtr = tail.node.nextPtr

		if tail == q.tail {
			if nextPtr.node == nil {
				newPtr := &linkedListNodePointer{
					node:  newNode,
					count: nextPtr.count + 1,
				}

				if atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&tail.node.nextPtr)),
					unsafe.Pointer(nextPtr),
					unsafe.Pointer(newPtr),
				) {
					break
				}

			} else {
				_ = atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
					unsafe.Pointer(tail),
					unsafe.Pointer(&linkedListNodePointer{
						node:  nextPtr.node,
						count: tail.count + 1,
					}),
				)
			}
		}
	}

	_ = atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
		unsafe.Pointer(tail),
		unsafe.Pointer(&linkedListNodePointer{
			node:  newNode,
			count: tail.count + 1,
		}),
	)

	return true
}

// deque on empty queue will return (nil, false)
func (q *ConcurrentLinkedListQueue) Dequeue() (value interface{}, ok bool) {
	fmt.Println("dequeuing")

	var (
		head *linkedListNodePointer
		tail *linkedListNodePointer
		next *linkedListNodePointer
	)

	for {
		head = q.head
		tail = q.tail
		next = head.node.nextPtr
		if head == q.head {
			if head.node == tail.node {
				if next.node == nil {
					return nil, false
				}
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(q.tail.node)),
					unsafe.Pointer(tail.node),
					unsafe.Pointer(&linkedListNodePointer{
						node:  next.node,
						count: tail.count + 1,
					}),
				)
			} else {
				value = next.node.value

				ok := atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&q.head)),
					unsafe.Pointer(head),
					unsafe.Pointer(&linkedListNodePointer{
						node:  next.node,
						count: head.count + 1,
					}),
				)

				if ok {
					break
				}
			}
		}
	}

	return value, true
}

func (q *ConcurrentLinkedListQueue) Len() int {
	tail := q.tail
	head := q.head
	return tail.count - head.count
}
