package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
	"os"
)

// SSTableWriter assumes that the keys are sorted.
// key and values exceeding maximum size will be rejected.
type SSTableWriter interface {
	Write(key types.KeyType, value types.ValueType, deleted bool, timestamp uint64) error
	MaxKeySize() uint32
	MaxValueSize() uint32
}

type SSTableReader interface {
}

type SSTableScanner interface {
	FindRecord(key types.KeyType) (value types.ValueType, deleted bool, ok bool)
}

type sstableBlockIndexWriter struct {
	FilePath string
	Storage  blockstore.BlockStorage
}

type (
	writer *sstableBlockIndexWriter
)

func NewIndexWriter(storage blockstore.BlockStorage, blockSize int, maxKeySize int) sstableBlockIndexWriter {
	return sstableBlockIndexWriter{
		Storage: storage,
	}
}

type sstableDataWriter struct {
	FilePath string
	File     *os.File
}

func (writer *sstableDataWriter) Write() (index blockstore.BlockOffSet,err error) {

}
