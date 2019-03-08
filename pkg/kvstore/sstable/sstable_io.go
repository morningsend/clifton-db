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
}

type SSTableReader interface {
	ReadNext() (key types.KeyType, value types.ValueType, deleted bool, err error)
	FindRecord(key types.KeyType) (value types.ValueType, deleted bool, ok bool, err error)
}
