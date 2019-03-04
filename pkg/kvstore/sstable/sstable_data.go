package sstable

import "github.com/zl14917/MastersProject/pkg/kvstore/types"

const DataFileMagic uint32 = 0x33333333

type SSTableDataFileHeader struct {
	Magic      uint32
	Flags      DataFileFlags
	BlockSize  uint32
	BlockCount uint32
}

type SSTableDataRecord struct {
	TimeStamp uint64
	Value     types.ValueType
}

type SSTableDataFile struct {
	Header       SSTableDataFileHeader
	RecordsCount uint32
}
