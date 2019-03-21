package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
)

// SSTableWriter assumes that the keys are sorted.
// key and values exceeding maximum size will be rejected.
type SSTableWriter interface {
	Write(key types.KeyType, value types.ValueType, deleted bool) error
	MaxKeySize() int
	MaxValueSize() int
	Commit() error
}

type SSTableReader interface {
	ReadNext() (key types.KeyType, value types.ValueType, deleted bool, err error)
	FindRecord(key types.KeyType) (value types.ValueType, deleted bool, ok bool, err error)
}


func NextMultipleOf4Uint(n uint) uint {
	return (n + 3) &^ uint(0x03)
}
