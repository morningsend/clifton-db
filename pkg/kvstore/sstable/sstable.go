package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
)

type SSTable struct {
}

type SSTableIndexFileHeader struct {
	FirstKey []byte
	LaskKey  []byte
	KeyCount uint32
}

type SSTableIndexFile struct {
	Header SSTableIndexFileHeader
}

type SSTableIndexEntry struct {
	Key    []byte
	OffSet uint64
}

type SSTableDataFileHeader struct {
}

type SSTableDataRecord struct {
	TimeStamp uint64
	Value     types.ValueType
}

type SSTableDataFile struct {
	Header       SSTableDataFileHeader
	RecordsCount uint32
}

func (f *SSTableDataFile) ReadRecordAt(offset uint64) SSTableDataRecord {

}
