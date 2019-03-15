package tables

type MemTableOps interface {
	Put(key []byte, value []byte) error
	Get(key []byte) (value [] byte, ok bool, err error)
	Remove(key []byte) (ok bool, err error)
}

type MemTable interface {
	MemTableOps

	KeyCountEstimate() uint
	Iterator() SortedKVIterator
}
