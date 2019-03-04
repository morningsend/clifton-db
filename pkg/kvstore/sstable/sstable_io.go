package sstable

import "github.com/zl14917/MastersProject/pkg/kvstore/types"

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
