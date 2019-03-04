package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
)

type EntryFlags uint32
type IndexFileFlags uint32
type DataFileFlags uint32

const (
	IndexFileMagic      uint32 = 0x32323232
	DataFileMagic       uint32 = 0x33333333
	HeaderUninitialized uint32 = 0x77777777
	BaseBlockSize       uint32 = 1024 * 4
)

const (
	KeyExists uint32 = (1 << iota)
	KeyDeleted
)

type SSTable struct {
	IndexFilePath string
	DataFilePath  string

	MaxKeySize   uint32
	MaxValueSize uint32
}

type SSTableOps interface {
	NewReader() SSTableReader
	NewWriter() SSTableWriter
}

type SSTableIndexFileHeader struct {
	Magic      uint32
	Flags      IndexFileFlags
	KeyCount   uint32
	BlockSize  uint32
	BlockCount uint32
	MaxKeySize uint32
}

type SSTableIndexFile struct {
	SSTableIndexFileHeader
	IndexBlocks []SSTableIndexBlock
}

type SSTableIndexBlock struct {
	KeyCount uint32
}

type SSTableIndexEntry struct {
	Flags          uint32
	KeyLen         uint32
	DataFileOffSet uint64
	SmallKey       [32]byte
	LargeKey       []byte
}

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

func NewSSTable() SSTable {
	return SSTable{
	}
}

func (s *SSTable) NewReader() SSTableReader {
	return nil
}

func (s *SSTable) NewWriter() SSTableWriter {
	return nil
}
