package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"os"
	"path"
)

type EntryFlags uint32
type DataFileFlags uint32

const (
	HeaderUninitialized uint32 = 0x77777777
	BaseBlockSize       uint32 = 1024 * 4
)

const (
	KeyExists uint32 = (1 << iota)
	KeyDeleted
)

type SSTable struct {
	IndexFilePath string
	DataFilePath  string

	MaxKeySize   int
	MaxValueSize int

	InMem            bool
	StorageBlockSize int
	dataStorage      blockstore.BlockStorage
	indexStorage     blockstore.BlockStorage
}

type SSTableOps interface {
	NewReader() SSTableReader
	NewWriter() SSTableWriter
	NewScanner() SSTableScanner
}

const indexFileName = "-index"
const dataFileName = "-data"

func UseInMemStore() {

}

type SSTableOpenOptions struct {
	Prefix       string
	MaxKeySize   int
	MaxValueSize int
	BlockSize    int
	KeyBlockSize int
	InMemStore   bool
}

func NewSSTable(dirPath string, options *SSTableOpenOptions) (*SSTable) {

	sstable := &SSTable{
		IndexFilePath:    path.Join(dirPath, options.Prefix+indexFileName),
		DataFilePath:     path.Join(dirPath, options.Prefix+dataFileName),
		InMem:            options.InMemStore,
		MaxKeySize:       options.MaxKeySize,
		MaxValueSize:     options.MaxValueSize,
		StorageBlockSize: options.BlockSize,
		indexStorage:     nil,
		dataStorage:      nil,
	}

	return sstable
}

func (s *SSTable) openStorage(useInMem bool, path string) (blockstore.BlockStorage, error) {
	var (
		store blockstore.BlockStorage
		err   error
	)
	if useInMem {
		store = blockstore.NewInMemBlockStorage(s.StorageBlockSize)
	} else {
		store, err = blockstore.NewBlockFile(
			path,
			blockstore.WithBlockSize(s.StorageBlockSize),
		)
	}

	if err != nil {
		return nil, err
	}
	return store, err
}
func (s *SSTable) createOrOpenIndexStorage(useInMem bool) (blockstore.BlockStorage, error) {

	return s.openStorage(useInMem, s.IndexFilePath)
}

func (s *SSTable) createOrOpenDataStorage(useInMem bool) (blockstore.BlockStorage, error) {
	return s.openStorage(useInMem, s.DataFilePath)
}

func LoadSSTableFrom(dirPath string, prefix string) *SSTable {
	sstable := &SSTable{
	}

	return sstable
}

func (s *SSTable) NewReader() SSTableReader {
	return nil
}

func (s *SSTable) NewScanner() SSTableScanner {
	return nil
}

func (s *SSTable) NewWriter() (SSTableWriter, error) {

	var (
		indexStorage blockstore.BlockStorage
		dataStorage  blockstore.BlockStorage
		err          error
	)
	if s.indexStorage == nil {
		indexStorage, err = s.createOrOpenIndexStorage(s.InMem)

		if err != nil {
			return nil, err
		}
	}

	dataStorage, err = s.createOrOpenDataStorage(s.InMem)

	if err != nil {
		_ = indexStorage.Close()

		if !s.InMem {
			_ = os.Remove(s.IndexFilePath)
		}
		return nil, err
	}

	s.indexStorage = indexStorage
	s.dataStorage = dataStorage

	return
}

func (s *SSTable) Close() error {
	return nil
}
