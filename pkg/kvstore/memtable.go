package kvstore

type MemTableOps interface {
	Put(key []byte, value []byte)
	Delete(key []byte, value []byte)
	Get(key []byte) (value [] byte, ok bool)
}

type MemTable interface {
	KeySizeEstimate() uint32
	BytesCount()
	Iterator() SortedKVIterator
}

func NewLockFreeMemTable() MemTable {
	return nil
}
