package kvstore

import (
	"context"
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
	"github.com/zl14917/MastersProject/pkg/kvstore/wal"
	"sync/atomic"
	"time"
	"unsafe"
)

type KVStoreReadOptions interface {
}

type KVStore interface {
	Get(keyType types.KeyType) (data types.ValueType, ok bool, err error)
	Put(key types.KeyType, data types.ValueType) (err error)
	Delete(key types.KeyType) (ok bool, err error)
	Exists(key types.KeyType) (ok bool, err error)
}

type KVStoreOpenOptions interface {
	Apply(options *KVStoreOptions)
}

type KVStoreOptions struct {
	WALSegmentSizeBytes int
	MaxKeySizeBytes     int
	MaxValueSizeBytes   int
	BlockSize           int
	IndexBlockSize      int
}

var defaultKVStoreOptions = KVStoreOptions{
	WALSegmentSizeBytes: 1024 * 1024 * 16,
	MaxKeySizeBytes:     1024 * 4,
	MaxValueSizeBytes:   1024 * 4,
}

type KVStoreMetadata struct {
	SStableLevel0 []string
	SStableLevel1 []string
	SStableLevel2 []string
	SStableLevel3 []string
}

type CliftonDBKVStore struct {
	fileTable FileTable
	memtable  MemTable
	wal       *wal.WAL

	backgroundCtx context.Context
	prevMemtable  MemTable
}

func (s *CliftonDBKVStore) bootstrapFromDir(dirPath string) error {
	return nil
}

func (s *CliftonDBKVStore) walCheckForRecovery() error {
	return nil
}

func (s *CliftonDBKVStore) rebuildMemTableFromWAL() error {
	return nil
}

func (s *CliftonDBKVStore) flushMemTable() error {

	newMemTable := NewLockFreeMemTable()
	s.prevMemtable = s.memtable

	for !atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&s.memtable)),
		unsafe.Pointer(&s.memtable),
		unsafe.Pointer(&newMemTable),
	) {

	}

	s.fileTable.BeginFlushing(s.prevMemtable, nil)

	return nil
}

func (s *CliftonDBKVStore) scheduleCompaction(deadline time.Duration) {

}

func FromDir(dirPath, options KVStoreOpenOptions) KVStore {
	store := &CliftonDBKVStore{
		memtable: nil,
		wal:      nil,
	}

	return store
}

func (CliftonDBKVStore) Get(keyType types.KeyType) (data types.ValueType, ok bool, err error) {
	return
}

func (CliftonDBKVStore) Put(key types.KeyType, data types.ValueType) (err error) {
	return
}

func (CliftonDBKVStore) Delete(key types.KeyType) (ok bool, err error) {
	return
}

func (CliftonDBKVStore) Exists(key types.KeyType) (ok bool, err error) {
	return
}
