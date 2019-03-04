package sstable

type SSTableWriter interface {
}

type SSTableReader interface {
	FindRecord(key []byte) (value []byte, deleted bool, ok bool)
}
