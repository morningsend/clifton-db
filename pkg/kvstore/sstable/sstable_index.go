package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"io"
	"os"
	"unsafe"
)

const (
	indexHeaderSize      = unsafe.Sizeof(SSTableIndexFileHeader{})
	blockHeaderSize      = unsafe.Sizeof(SSTableIndexBlock{})
	indexEntryHeaderSize = unsafe.Sizeof(SSTableIndexEntry{}.Flags) +
		unsafe.Sizeof(SSTableIndexEntry{}.KeyLen) +
		unsafe.Sizeof(SSTableIndexEntry{}.DataFileOffSet)
	IndexFileMagic uint32 = 0x32323232
)

const (
	SSTableIndexKeyInsert IndexKeyFlags = 1 << iota
	SSTableIndexKeyDelete
)

var InvalidHeaderMagicErr = errors.New("first 32-bit magic of file is wrong")

type IndexKeyFlags uint32
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
	blockBuffer *bytes.Buffer
	storage     blockstore.BlockStorage
}

type SSTableIndexBlock struct {
	KeyCount uint32
}

type SSTableIndexEntry struct {
	Flags          IndexKeyFlags
	KeyLen         uint32
	DataFileOffSet uint64
	LargeKey       []byte
}

func (header *SSTableIndexFileHeader) Marshall(writer io.Writer) error {
	var (
		smallBuffer  [4]byte
		uint32buffer = smallBuffer[0:4]
	)

	binary.BigEndian.PutUint32(uint32buffer, header.Magic)

	_, err := writer.Write(uint32buffer)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buffer, uint32(header.Flags))

	_, err = writer.Write(uint32buffer)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(uint32buffer, header.KeyCount)

	_, err = writer.Write(uint32buffer)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(uint32buffer, header.BlockSize)

	_, err = writer.Write(uint32buffer)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(uint32buffer, header.BlockCount)

	_, err = writer.Write(uint32buffer)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buffer, header.MaxKeySize)

	_, err = writer.Write(uint32buffer)
	if err != nil {
		return err
	}

	return nil
}

func (header *SSTableIndexFileHeader) UnMarshall(reader io.Reader) error {
	var (
		smallbuffer [4]byte
		uint32buf   = smallbuffer[0:4]
	)

	_, err := reader.Read(uint32buf)
	if err != nil {
		return err
	}

	header.Magic = binary.BigEndian.Uint32(uint32buf)
	if header.Magic != IndexFileMagic {
		return InvalidHeaderMagicErr
	}

	_, err = reader.Read(uint32buf)
	if err != nil {
		return err
	}
	header.Flags = IndexFileFlags(binary.BigEndian.Uint32(uint32buf))

	_, err = reader.Read(uint32buf)
	if err != nil {
		return err
	}
	header.KeyCount = binary.BigEndian.Uint32(uint32buf)

	_, err = reader.Read(uint32buf)
	if err != nil {
		return err
	}
	header.BlockSize = binary.BigEndian.Uint32(uint32buf)

	_, err = reader.Read(uint32buf)
	if err != nil {
		return err
	}
	header.BlockCount = binary.BigEndian.Uint32(uint32buf)

	_, err = reader.Read(uint32buf)
	if err != nil {
		return err
	}
	header.MaxKeySize = binary.BigEndian.Uint32(uint32buf)

	return nil
}

func (e *SSTableIndexEntry) Marshall(buffer io.Writer) (nbytes int, err error) {
	if e.KeyLen != uint32(len(e.LargeKey)) {
		return 0, fmt.Errorf(
			"Entry KeyLen is %d does not equal length of key %d",
			e.KeyLen,
			len(e.LargeKey),
		)
	}

	nbytes = 0

	var (
		smallBytesBuffer [8]byte
		uint32buffer     = smallBytesBuffer[0:4]
		uint64buffer     = smallBytesBuffer[0:8]
	)

	binary.BigEndian.PutUint32(uint32buffer, uint32(e.Flags))
	_, err = buffer.Write(uint32buffer)
	if err != nil {
		return
	}

	binary.BigEndian.PutUint32(uint32buffer, e.KeyLen)
	_, err = buffer.Write(uint32buffer)

	if err != nil {
		return
	}

	binary.BigEndian.PutUint64(uint64buffer, e.DataFileOffSet)
	_, err = buffer.Write(uint64buffer)

	if err != nil {
		return
	}

	nbytes += int(indexEntryHeaderSize)

	if len(e.LargeKey) < 1 {
		return
	}

	keyLengthPadded := NextMultipleOf4Uint(uint(e.KeyLen))
	keyPadToSize := keyLengthPadded - uint(e.KeyLen)
	binary.BigEndian.PutUint32(uint32buffer, 0)

	_, err = buffer.Write(e.LargeKey)
	if err != nil {
		return
	}

	_, err = buffer.Write(uint32buffer[0:keyPadToSize])
	nbytes += int(keyLengthPadded)

	return nbytes, nil
}

func (e *SSTableIndexEntry) UnMarshall(buffer *bytes.Buffer) error {
	var (
		err error

		smallBytesBuffer [8]byte
		uint32buffer     = smallBytesBuffer[0:4]
		uint64buffer     = smallBytesBuffer[0:8]
	)

	_, err = buffer.Read(uint32buffer)
	if err != nil {
		return err
	}

	e.Flags = IndexKeyFlags(binary.BigEndian.Uint32(uint32buffer))

	_, err = buffer.Read(uint32buffer)
	if err != nil {
		return err
	}
	e.KeyLen = binary.BigEndian.Uint32(uint32buffer)

	_, err = buffer.Read(uint64buffer)
	if err != nil {
		return err
	}

	e.DataFileOffSet = binary.BigEndian.Uint64(uint64buffer)
	paddedKeySize := NextMultipleOf4Uint(uint(e.KeyLen))
	e.LargeKey = make([]byte, paddedKeySize, paddedKeySize)
	n, err := buffer.Read(e.LargeKey)

	if err != nil {
		return err
	}

	if uint(n) != paddedKeySize {
		return fmt.Errorf("expect KeyLen %d but only read %d", e.KeyLen, n)
	}

	return nil
}

func validateBlockSize(blockSize uint32) bool {
	return blockSize%BaseBlockSize == 0
}

func NewSSTableIndexFile(keyBlockSize uint32) (SSTableIndexFile, error) {

	if !validateBlockSize(keyBlockSize) {
		return SSTableIndexFile{},
			fmt.Errorf(
				"key block size must be a multiple of %d, got %d",
				BaseBlockSize,
				keyBlockSize)
	}

	indexFile := SSTableIndexFile{

		SSTableIndexFileHeader: SSTableIndexFileHeader{
			Magic:      IndexFileMagic,
			Flags:      IndexFileFlags(HeaderUninitialized),
			MaxKeySize: uint32(MaxKeySizeFitInBlocK(int(keyBlockSize))),
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

func MaxKeySizeFitInBlocK(blockSize int) int {
	availableBlockBytes := blockSize - int(blockHeaderSize)
	availableEntryBytes := availableBlockBytes - int(indexEntryHeaderSize)

	return availableEntryBytes
}

func SizeOfIndexEntry(entry *SSTableIndexEntry) int {
	return int(indexEntryHeaderSize) + int(entry.KeyLen)
}

func MarshalledSizeOfIndexEntry(entry *SSTableIndexEntry) int {
	return int(indexEntryHeaderSize) + int(NextMultipleOf4Uint(uint(entry.KeyLen)))
}
