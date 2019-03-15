package tables

import "github.com/zl14917/MastersProject/pkg/kvstore/types"

type SortedKVIterator interface {
	Next() bool
	Current() (key types.KeyType, value types.ValueType, flags int)
}


