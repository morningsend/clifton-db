package blockstore

import (
	"bytes"
	"errors"
	"io"
)

var SizeExceedBlockSize = errors.New("size exceeded")

const BaseBlockSize = 4 * 1024

// Random access to blocks in a file.
// IO is aligned.
// Returns number of bytes read or error
type BlockReader interface {
	ReadBlock(index uint, buffer *bytes.Buffer) (n int, err error)
}

// Random write to block in a file.
// Only write data up to block size, or up to size of buffer,
// which ever is smaller
// IO is aligned.
// Returns number of bytes written or error
type BlockWriter interface {
	WriteBlock(index uint, buffer *bytes.Buffer) (n int, err error)
}

type BlockAllocator interface {
	Allocate(nblocks int) (nAllocated int, err error)
}

type Position struct {
	Block  int
	Offset int
}

var UninitializedPosition = Position{
	-1, -1,
}

type BlockStorage interface {
	blockStorage()

	BlockReader
	BlockWriter
	BlockAllocator

	io.Reader
	io.Writer
	io.Closer

	NumBlocks() int
	BlockSize() int

	WritePosition() Position
	ReadPosition() Position

	Sync() error
	Flush() error
}
