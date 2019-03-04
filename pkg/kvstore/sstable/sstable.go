package sstable

import (
	"os"
	"path"
)

type EntryFlags uint32
type DataFileFlags uint32

const (
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

const indexFileName = "-index"
const dataFileName = "-data"

func NewSSTable(dirPath string, prefix string, maxKeySize uint32, maxValueSize uint32, keyBlockSize uint32) (*SSTable, error) {
	var err error
	sstable := &SSTable{
		IndexFilePath: path.Join(dirPath, prefix+indexFileName),
		DataFilePath:  path.Join(dirPath, prefix+dataFileName),
		MaxKeySize:    maxKeySize,
		MaxValueSize:  maxValueSize,
	}

	err = sstable.CreateIndexFile()
	if err != nil {
		return nil, err
	}

	err = sstable.CreateDataFile()

	if err != nil {
		_ = os.Remove(sstable.IndexFilePath)
		return nil, err
	}

	return sstable, nil
}

func (s *SSTable) CreateIndexFile() error {
	return s.createFile(s.IndexFilePath)
}

func (s *SSTable) createFile(path string) error {
	file, err := os.OpenFile(
		s.IndexFilePath,
		os.O_CREATE|os.O_EXCL|os.O_RDWR,
		0644,
	)

	if err != nil {
		return err
	}

	err = file.Close()
	return err
}

func (s *SSTable) CreateDataFile() error {
	return s.createFile(s.DataFilePath)
}

func LoadSSTableFrom(dirPath string, prefix string) *SSTable {
	sstable := &SSTable{

	}

	return sstable
}

func (s *SSTable) NewReader() SSTableReader {
	return nil
}

func (s *SSTable) NewScanner() SSTableScanner {
	return nil
}

func (s *SSTable) NewWriter() SSTableWriter {
	return nil
}

func (s *SSTable) Close() error {
	return nil
}
