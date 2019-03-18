package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
	"unsafe"
)

const (
	DataFileMagic uint32 = 0x33333333

	SSTableDataFileHeaderSize   = unsafe.Sizeof(SSTableDataFileHeader{})
	SSTableDataRecordHeaderSize = unsafe.Sizeof(SSTableDataRecord{}.Value)
)

type SSTableDataFileHeader struct {
	Magic      uint32
	Flags      DataFileFlags
	BlockSize  uint32
	BlockCount uint32
}

type SSTableDataRecord struct {
	ValueLen uint32
	Value    types.ValueType
}

type SSTableDataFile struct {
	Header       SSTableDataFileHeader
	RecordsCount uint32
}

func MaxValueSizeFitInBlock(blockSize int) int {
	return blockSize - int(SSTableDataRecordHeaderSize)
}
