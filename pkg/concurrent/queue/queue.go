package queue

type Queue interface {
	Enqueue(value interface{}) (ok bool)
	Dequeue() (value interface{}, ok bool)
	Len() int
}
