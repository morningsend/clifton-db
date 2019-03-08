package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
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

type sstableBlockIndexFileWriter struct {
	FilePath string
	Storage  blockstore.BlockStorage
}

func newIndexWriter(storage blockstore.BlockStorage, blockSize int, maxKeySize int) sstableBlockIndexFileWriter {
	return sstableBlockIndexFileWriter{
		Storage: storage,
	}
}

type sstableDataFileWriter struct {
	Storage   blockstore.BlockStorage
	BlockSize int
}

func (writer *sstableDataFileWriter) Write(value []byte) (index blockstore.Position, err error) {
	return blockstore.Position{}, nil
}

func newSStableDataWriter(storage blockstore.BlockStorage, blockSize int) sstableDataFileWriter {
	return sstableDataFileWriter{
		Storage:   storage,
		BlockSize: blockSize,
	}
}
