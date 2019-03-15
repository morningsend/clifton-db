package kvstore

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
	"github.com/zl14917/MastersProject/pkg/kvstore/wal"
	"github.com/zl14917/MastersProject/pkg/logger"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	sstablePath  = "sstables/"
	walPath      = "wal/"
	logFileName  = "kvstore-%d.log"
	logPrefix    = "[kvstore-%d]"
	lockFileName = "store.lock.file"
)

type KVStoreLockFileData struct {
	PartitionId    uint32 `yaml:"partition-id"`
	WALCommitIndex uint64 `yaml:"wal-commit-index"`
	WALApplyIndex  uint64 `yaml:"wal-apply-index"`
}

var defaultKVStoreLockFileData = KVStoreLockFileData{
	PartitionId:    0,
	WALCommitIndex: 0,
	WALApplyIndex:  0,
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
	DataBlockSize       int
	IndexBlockSize      int
}

var defaultKVStoreOptions = KVStoreOptions{
	WALSegmentSizeBytes: 1024 * 1024 * 16,
	DataBlockSize:       1024 * 16,
	IndexBlockSize:      1024 * 4,
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

	logger        logger.Logger
	backgroundCtx context.Context
	prevMemtable  MemTable

	KVStoreRoot         string
	SSTablesRoot        string
	WALRoot             string
	KVStoreLockFilePath string
}

func NewCliftonDBKVStore(dirPath string, logPath string) (*CliftonDBKVStore, error) {
	var err error

	walRootPath := path.Join(dirPath, walPath)

	store := &CliftonDBKVStore{
		fileTable:    nil,
		memtable:     NewMapMemTable(),
		wal:          wal.NewWAL(walRootPath),
		KVStoreRoot:  dirPath,
		SSTablesRoot: path.Join(dirPath, sstablePath),

		WALRoot:             walRootPath,
		KVStoreLockFilePath: path.Join(dirPath, walRootPath, lockFileName),

		logger: nil,
	}

	data, err := store.ReadLockFile()

	storeLogFilePath := path.Join(logPath, fmt.Sprintf(logFileName, data.PartitionId))

	if err != nil {
		return nil, err
	}
	err = store.EnsureDirsExist()

	if err != nil {
		return nil, err
	}

	store.logger, err = logger.NewFileLogger(
		storeLogFilePath,
		fmt.Sprintf(logPrefix, data.PartitionId),
		log.LstdFlags,
	)

	if err != nil {
		store.logger = log.New(os.Stdout, fmt.Sprintf(logPrefix, data.PartitionId), log.LstdFlags)
	}
	err = store.walCheckForRecovery()

	if err != nil {
		return nil, err
	}
	err = store.WriteLockFile(data)

	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *CliftonDBKVStore) ReadLockFile() (data KVStoreLockFileData, err error) {
	file, err := os.OpenFile(s.KVStoreLockFilePath, os.O_RDONLY, 0644)

	if err != nil {
		return defaultKVStoreLockFileData, err
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&data)
	if err != nil {
		return defaultKVStoreLockFileData, err
	}
	return
}

func (s *CliftonDBKVStore) WriteLockFile(data KVStoreLockFileData) error {
	file, err := os.OpenFile(s.KVStoreLockFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	err = yaml.NewEncoder(file).Encode(data)
	if err != nil {
		_ = file.Close()
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *CliftonDBKVStore) EnsureDirsExist() error {
	err := os.MkdirAll(s.WALRoot, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.SSTablesRoot, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func BootstrapFromDir(dirPath string) *KVStore {
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

	err := s.fileTable.BeginFlushing(s.prevMemtable)

	if err != nil {
		return err
	}

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
