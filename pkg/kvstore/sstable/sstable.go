package sstable

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
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

var KeyTooLarge = errors.New("SSTable key too large to be stored")
var ValueTooLarge = errors.New("SSTable value too large to be stored")

type SSTable struct {
	IndexFilePath string
	DataFilePath  string

	MaxKeySize   int
	MaxValueSize int

	InMem            bool
	StorageBlockSize int

	dataStorage  blockstore.BlockStorage
	indexStorage blockstore.BlockStorage

	CreatedTimestamp uint64

	loadExisting bool
}

type SSTableOps interface {
	NewReader() SSTableReader
	NewWriter() SSTableWriter
}

const indexFileName = "-index"
const dataFileName = "-data"

type SSTableOpenOptions struct {
	Prefix       string
	MaxKeySize   int
	MaxValueSize int
	BlockSize    int
	KeyBlockSize int
	InMemStore   bool
	LoadExisting bool
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
		loadExisting:     options.LoadExisting,
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

		return store, nil
	}

	if s.loadExisting {
		store, err = blockstore.OpenBlockFile(
			path,
			blockstore.WithBlockSize(s.StorageBlockSize),
		)
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

func LoadSSTableFrom(dirPath string, options *SSTableOpenOptions) *SSTable {
	options.LoadExisting = true
	return NewSSTable(dirPath, options)
}

func (s *SSTable) NewReader() (SSTableReader, error) {
	reader := &sstableReaderStruct{

	}
	return reader, nil
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
		s.indexStorage = indexStorage
	}

	if s.dataStorage == nil {
		dataStorage, err = s.createOrOpenDataStorage(s.InMem)

		if err != nil {
			_ = indexStorage.Close()

			if !s.InMem {
				_ = os.Remove(s.IndexFilePath)
			}
			return nil, err
		}
		s.dataStorage = dataStorage
	}

	indexWriteBuffer := bytes.NewBuffer(nil)
	indexWriteBuffer.Grow(s.StorageBlockSize)

	writer := &sstableWriterStruct{
		sstableDataWriter: sstableDataWriter{
			Storage:      dataStorage,
			BlockSize:    s.StorageBlockSize,
			MaxValueSize: s.MaxValueSize,
		},
		sstableBlockIndexWriter: sstableBlockIndexWriter{
			Storage:    s.indexStorage,
			BlockSize:  s.StorageBlockSize,
			buffer:     indexWriteBuffer,
			MaxKeySize: s.MaxKeySize,
		},
	}

	return writer, nil
}

func (s *SSTable) Close() error {
	var errIndex, errData error
	if s.indexStorage != nil {
		errIndex = s.indexStorage.Close()
	}
	if s.dataStorage != nil {
		errData = s.dataStorage.Close()
	}

	if errIndex != nil || errData != nil {
		return fmt.Errorf("error closing storages:", errIndex, errData)
	}

	return nil
}

type sstableBlockIndexWriter struct {
	FilePath   string
	Storage    blockstore.BlockStorage
	BlockSize  int
	buffer     *bytes.Buffer
	MaxKeySize int
}

func (writer *sstableBlockIndexWriter) writeIndex(key types.KeyType, deleted bool, position blockstore.Position) error {
	return nil
}

func newIndexWriter(storage blockstore.BlockStorage, blockSize int, maxKeySize int) sstableBlockIndexWriter {
	return sstableBlockIndexWriter{
		Storage: storage,
	}
}

type sstableDataWriter struct {
	Storage      blockstore.BlockStorage
	BlockSize    int
	MaxValueSize int
}

func (writer *sstableDataWriter) writeValue(value types.ValueType) (index blockstore.Position, err error) {
	return blockstore.Position{}, nil
}
func newSStableDataWriter(storage blockstore.BlockStorage, blockSize int) sstableDataWriter {
	return sstableDataWriter{
		Storage:   storage,
		BlockSize: blockSize,
	}
}

type sstableWriterStruct struct {
	sstableDataWriter
	sstableBlockIndexWriter
}

func (w *sstableWriterStruct) MaxKeySize() int {
	return w.sstableBlockIndexWriter.MaxKeySize
}

func (w *sstableWriterStruct) MaxValueSize() int {
	return w.sstableDataWriter.MaxValueSize
}

func (w *sstableWriterStruct) Write(key types.KeyType, value types.ValueType, deleted bool) error {
	var (
		keyLen   = len(key)
		valueLen = len(value)
	)

	if w.sstableBlockIndexWriter.MaxKeySize < keyLen {
		return KeyTooLarge
	}

	if w.sstableDataWriter.MaxValueSize < valueLen {
		return ValueTooLarge
	}

	position, err := w.writeValue(value)

	if err != nil {
		return err
	}

	err = w.writeIndex(key, deleted, position)
	return err
}

type sstableReaderStruct struct {
	IndexStorage blockstore.BlockStorage
	DataStorage  blockstore.BlockStorage
}

func (r *sstableReaderStruct) ReadNext() (key types.KeyType, value types.ValueType, deleted bool, err error) {
	return
}

func (r *sstableReaderStruct) FindRecord(key types.KeyType) (value types.ValueType, deleted bool, ok bool, err error) {
	return
}
