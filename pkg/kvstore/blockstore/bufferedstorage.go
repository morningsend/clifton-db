package blockstore

import (
	"bytes"
	"os"
)

const (
	FsBaseBlockSize uint = 4 * 1024
)

// A class to provide block buffering for io.
// Data is aligned in blocks
type BufferedBlockStorage struct {
	FilePath string
	fileIO   *os.File

	blockCount uint
	autoSync   bool
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

func NewBlockFile(path string, storageOptions ...BufferedBlockStorageOption) (BlockStorage, error) {
	var options = defaultOptions
	applyOptions(&options, storageOptions...)

	storage := &BufferedBlockStorage{

	}

	return storage, nil
}

func OpenBlockFile(path string, storageOptions ...BufferedBlockStorageOption) (BlockStorage, error) {
	var options = defaultOptions
	applyOptions(&options, storageOptions...)

	storage := &BufferedBlockStorage{

	}

	return storage, nil
}

func (s *BufferedBlockStorage) Read(data []byte) (n int, err error) {
	return
}

func (s *BufferedBlockStorage) Write(data []byte) (n int, err error) {
	return
}

func (s *BufferedBlockStorage) Sync() error {
	return nil
}

func (s *BufferedBlockStorage) Close() error {
	return nil
}
func (s *BufferedBlockStorage) Allocate(nblocks int) (nAllocated int, err error) {
	return
}

func (s *BufferedBlockStorage) Flush() error {
	return nil
}

func (s *BufferedBlockStorage) NumBlocks() int {
	return 0
}

func (s *BufferedBlockStorage) BlockSize() int {
	return 0
}

func (s *BufferedBlockStorage) WriteBlock(index uint, buffer *bytes.Buffer) (n int, err error) {
	return
}

func (s *BufferedBlockStorage) ReadBlock(index uint, buffer *bytes.Buffer) (n int, err error) {
	return
}
