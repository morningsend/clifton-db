package tables

import (
	"github.com/zl14917/MastersProject/concurrent/maps"
	"github.com/zl14917/MastersProject/kvstore/types"
)

type SortedKVIterator interface {
	Next() bool
	Current() (key types.KeyType, value maps.Value)
}
