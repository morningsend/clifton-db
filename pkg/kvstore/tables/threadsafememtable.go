package tables

import "github.com/zl14917/MastersProject/pkg/kvstore/maps/concurrent"

type ThreadSafeMapMemTable struct {
	concurrent.ThreadsafeMap
	MaxKeySize   int
	MaxValueSize int
}

func (ThreadSafeMapMemTable) Put(key []byte, value []byte) error {
	return nil
}

func (ThreadSafeMapMemTable) Get(key []byte) (value [] byte, ok bool, err error) {
	return
}

func (ThreadSafeMapMemTable) Remove(key []byte) (ok bool, err error) {
	return
}

func (m *ThreadSafeMapMemTable) KeyCountEstimate() uint {
	return uint(m.Len())
}

func (ThreadSafeMapMemTable) Iterator() SortedKVIterator {
	return nil
}

func NewMapMemTable(maxKeySize int, maxValueSize int) MemTable {
	return &ThreadSafeMapMemTable{
		MaxKeySize:   maxKeySize,
		MaxValueSize: maxValueSize,
	}
}
