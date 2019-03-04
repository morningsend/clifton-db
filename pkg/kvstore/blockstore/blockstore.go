package blockstore

import (
	"bytes"
	"errors"
	"io"
)

const (
	BaseBlockSize uint = 4 * 1024
)

var SizeExceeded = errors.New("size exceeded")

// Random access to blocks in a file.
// IO is aligned.
type BlockReader interface {
	ReadBlock(index uint, buffer *bytes.Buffer) (n int, err error)
}

// Random write to block in a file.
// IO is aligned.
type BlockWriter interface {
	WriteBlock(index uint, buffer *bytes.Buffer) (n int, err error)
}

type BlockStorage interface {
	blockStorage()

	BlockReader
	BlockWriter

	io.Reader
	io.Writer
	io.Closer

	NumBlocks() int
	BlockSize() int

	Sync() error
	Flush() error
}
