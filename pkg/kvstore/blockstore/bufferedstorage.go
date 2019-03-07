package blockstore

import (
	"bytes"
	"io"
	"os"
)

const (
	FsBaseBlockSize uint = 4 * 1024
)

// A class to provide block buffering for io.
// Data is aligned in blocks
type BufferedBlockStorage struct {
	FilePath string
	file     *os.File

	blockSize int
	blockLen  int
	autoSync  bool
	syncMode  bool

	seqReadBlock   int
	seqReadOffset  int
	seqWriteBlock  int
	seqWriteOffset int

	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer

	currentReadBufferBlock  int
	currentWriteBufferBlock int
}

type BufferedBlockStorageOption func(options *BufferedBlockStorageOptions)

type BufferedBlockStorageOptions struct {
	AutoSync   bool
	SyncFileIO bool
	BlockSize  int
}

var defaultOptions = BufferedBlockStorageOptions{
	AutoSync:   false,
	SyncFileIO: false,
	BlockSize:  int(FsBaseBlockSize),
}

func WithBlockSize(value int) BufferedBlockStorageOption {
	return func(opts *BufferedBlockStorageOptions) {
		opts.BlockSize = value
	}
}

func WithSyncFileIO() BufferedBlockStorageOption {
	return func(opts *BufferedBlockStorageOptions) {
		opts.SyncFileIO = true
	}
}

func WithAutoSync() BufferedBlockStorageOption {
	return func(opts *BufferedBlockStorageOptions) {
		opts.AutoSync = true
	}
}

func (b *BufferedBlockStorage) blockStorage() {}

func applyOptions(opts *BufferedBlockStorageOptions, storageOptions ...BufferedBlockStorageOption) {
	for _, p := range storageOptions {
		p(opts)
	}
}

func newBufferedStorageWithOptions(path string, storageOptions ...BufferedBlockStorageOption) *BufferedBlockStorage {
	var options = defaultOptions
	applyOptions(&options, storageOptions...)
	storage := &BufferedBlockStorage{
		seqReadOffset:  0,
		seqReadBlock:   0,
		seqWriteBlock:  0,
		seqWriteOffset: 0,

		FilePath: path,

		autoSync: !options.SyncFileIO && options.AutoSync,
		syncMode: options.SyncFileIO,

		blockLen: 0,

		readBuffer:              bytes.NewBuffer(nil),
		writeBuffer:             bytes.NewBuffer(nil),
		currentReadBufferBlock:  0,
		currentWriteBufferBlock: 0,
	}

	storage.readBuffer.Grow(storage.blockSize)
	storage.writeBuffer.Grow(storage.blockSize)

	return storage
}

func NewBlockFile(path string, storageOptions ...BufferedBlockStorageOption) (BlockStorage, error) {
	var err error
	var flags int = os.O_RDWR | os.O_CREATE | os.O_EXCL
	var options = defaultOptions

	applyOptions(&options, storageOptions...)

	storage := newBufferedStorageWithOptions(path, storageOptions...)

	if storage.syncMode {
		flags |= os.O_SYNC
	}

	storage.file, err = os.OpenFile(
		storage.FilePath,
		flags,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return storage, nil
}

func OpenBlockFile(path string, storageOptions ...BufferedBlockStorageOption) (BlockStorage, error) {
	var err error
	var flags int = os.O_RDWR

	storage := newBufferedStorageWithOptions(path, storageOptions...)

	if storage.syncMode {
		flags |= os.O_SYNC
	}

	storage.file, err = os.OpenFile(
		storage.FilePath,
		flags,
		0644,
	)

	fileInfo, err := storage.file.Stat()
	storage.blockLen = int(fileInfo.Size() / int64(storage.blockSize))

	if fileInfo.Size()%int64(storage.blockSize) > 0 {
		storage.blockLen += 1
	}

	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *BufferedBlockStorage) flushReadBuffer() error {

	clearBuffer(s.readBuffer, 0x0, s.blockSize)
	s.readBuffer.Reset()
	bufMem := s.readBuffer.Bytes()

	// Read in a loop
	_, err := s.file.Read(bufMem[0:s.blockSize])

	if err != nil {
		return err
	}

	return nil
}

func (s *BufferedBlockStorage) Read(data []byte) (n int, err error) {
	bytesToRead := len(data)

	if bytesToRead > s.blockSize {
		return 0, SizeExceedBlockSize
	}

	remaining := s.blockSize - s.seqReadOffset

	// if reading bytes larger than space remaining in current block,
	// read the next block from disk into buffer.
	if bytesToRead > remaining {
		s.seqReadOffset = 0
		s.seqReadBlock++
		remaining = s.blockSize

		s.flushReadBuffer()
	}

	if s.blockLen <= s.seqReadBlock {
		return 0, io.EOF
	}

	return
}

func (s *BufferedBlockStorage) Write(data []byte) (n int, err error) {
	bytesToWrite := len(data)

	if bytesToWrite > s.blockSize {
		return 0, SizeExceedBlockSize
	}

	remaining := s.blockSize - s.seqWriteBlock

	// if data larger than space remaining in current block,
	// flush buffer to disk, increase cursor position,
	// and reset buffer.
	if bytesToWrite > remaining {
		s.seqWriteBlock++
		s.seqWriteOffset = 0
		err = s.Flush()
		if err != nil {
			return 0, err
		}

		clearBuffer(s.writeBuffer, 0x0, s.blockLen)
		s.writeBuffer.Reset()

		s.currentWriteBufferBlock = s.seqWriteBlock
	}

	// if don't have enough space in file, we grow as needed.
	if s.blockLen < s.seqWriteBlock+1 {
		oldLen := s.blockLen
		_, err = s.Allocate(s.seqWriteBlock + 10)
		if err != nil {
			s.blockLen = oldLen
			return 0, err
		}
	}

	s.writeBuffer.Write(data)
	s.seqWriteOffset += bytesToWrite

	return bytesToWrite, nil
}

func (s *BufferedBlockStorage) Sync() error {
	return s.file.Sync()
}

func (s *BufferedBlockStorage) Close() error {
	return s.file.Close()
}
func (s *BufferedBlockStorage) Allocate(nblocks int) (nAllocated int, err error) {
	if s.blockLen >= nblocks {
		return
	}

	growth := nblocks - s.blockLen
	size := int64(nblocks) * int64(s.blockLen)
	err = s.file.Truncate(size)

	if err != nil {
		return
	}
	s.blockLen = nblocks

	return growth, nil
}

func (s *BufferedBlockStorage) Flush() error {
	offset := int64(s.currentWriteBufferBlock) * int64(s.blockSize)
	_, err := s.file.WriteAt(s.writeBuffer.Bytes()[:s.blockSize], offset)
	return err
}

func (s *BufferedBlockStorage) NumBlocks() int {
	return s.blockLen
}

func (s *BufferedBlockStorage) BlockSize() int {
	return s.blockSize
}

func (s *BufferedBlockStorage) WriteBlock(index uint, buffer *bytes.Buffer) (n int, err error) {
	return
}

func (s *BufferedBlockStorage) ReadBlock(index uint, buffer *bytes.Buffer) (n int, err error) {
	return
}

func (s *BufferedBlockStorage) ReadPosition() Position {
	return Position{
		Block:  s.seqReadBlock,
		Offset: s.seqReadOffset,
	}
}

func (s *BufferedBlockStorage) WritePosition() Position {
	return Position{
		Block:  s.seqWriteBlock,
		Offset: s.seqWriteOffset,
	}
}
