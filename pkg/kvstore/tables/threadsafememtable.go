package tables

import "github.com/zl14917/MastersProject/pkg/kvstore/maps/concurrent"

type ThreadSafeMapMemTable struct {
	concurrent.ThreadsafeMap
	MaxKeySize   int
	MaxValueSize int
}

func (m *ThreadSafeMapMemTable) Exists(key []byte) (ok bool, err error) {
	panic("implement me")
}

func (m *ThreadSafeMapMemTable) Put(key []byte, value []byte) error {
	return nil
}

func (m *ThreadSafeMapMemTable) Get(key []byte) (value [] byte, ok bool, err error) {
	return
}

func (m *ThreadSafeMapMemTable) Remove(key []byte) (ok bool, err error) {
	return
}

func (m *ThreadSafeMapMemTable) KeyCountEstimate() uint {
	return uint(m.Len())
}

func (m *ThreadSafeMapMemTable) Iterator() SortedKVIterator {
	iter := m.ThreadsafeMap.Iterator()

}

func NewMapMemTable(maxKeySize int, maxValueSize int) MemTable {
	return &ThreadSafeMapMemTable{
		MaxKeySize:   maxKeySize,
		MaxValueSize: maxValueSize,
	}
}
