package maps

type Key string
type Value []byte

type MapIterator interface {
	Next() bool
	Current() (key Key, value Value)
}

type Map interface {
	Get(key Key) (value Value, ok bool)
	Put(key Key, value Value) (err error)
	Contains(key Key) (ok bool)
	Iterator() MapIterator
}
