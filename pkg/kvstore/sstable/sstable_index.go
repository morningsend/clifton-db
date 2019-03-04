package sstable

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unsafe"
)

var indexHeaderSize = unsafe.Sizeof(SSTableIndexFileHeader{})

const IndexFileMagic uint32 = 0x32323232

type IndexFileFlags uint32

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

	Path        string
	preloaded   bool
	fileIO      *os.File
	blockBuffer *bytes.Buffer
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

func (header *SSTableIndexFileHeader) Marshall(writer io.Writer) {

}

func (header *SSTableIndexFileHeader) UnMarshall(reader io.Reader) {

}

func validateBlockSize(blockSize uint32) bool {
	return blockSize%BaseBlockSize == 0
}

func NewSSTableIndexFile(path string, maxKeySize uint32, keyBlockSize uint32) (SSTableIndexFile, error) {

	if !validateBlockSize(keyBlockSize) {
		return SSTableIndexFile{},
			fmt.Errorf(
				"key block size must be a multiple of %d, got %d",
				BaseBlockSize,
				keyBlockSize)
	}

	indexFile := SSTableIndexFile{
		Path: path,

		SSTableIndexFileHeader: SSTableIndexFileHeader{
			Magic:      IndexFileMagic,
			Flags:      IndexFileFlags(HeaderUninitialized),
			MaxKeySize: maxKeySize,
			BlockSize:  keyBlockSize,
			KeyCount:   HeaderUninitialized,
			BlockCount: HeaderUninitialized,
		},

		IndexBlocks: make([]SSTableIndexBlock, 0, 8),

		blockBuffer: bytes.NewBuffer(nil),
	}

	err := indexFile.CreateFile()

	return indexFile, err
}

func LoadSSTableIndexFile(path string) (SSTableIndexFile, error) {
	return SSTableIndexFile{}, nil
}

func (index *SSTableIndexFile) Preload() error {
	return nil
}

// create file opens file and writes header
func (index *SSTableIndexFile) CreateFile() error {
	file, err := os.Create(index.Path)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte{})
	return err
}

func (index *SSTableIndexFile) WriteEntry() {

}
