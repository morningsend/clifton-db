package blockstore

import "bytes"

type inMemBlock struct {
	data []byte
}

func newInMemBlock(size int) inMemBlock {
	return inMemBlock{
		data: make([]byte, size, size),
	}
}

// Im Memory block storage is uesd for testing only
type InMemBlockStorage struct {
	blocks   []inMemBlock
	blockLen int

	blockSize int

	seqReadBlock  int
	seqReadOffset int

	seqWriteBlock  int
	seqWriteOffset int
}

func NewInMemBlockStorage(blockSize int) BlockStorage {
	return &InMemBlockStorage{
		blockSize:      blockSize,
		blockLen:       -1,
		blocks:         make([]inMemBlock, 0, 1),
		seqWriteBlock:  blockSize,
		seqWriteOffset: blockSize,
		seqReadBlock:   0,
		seqReadOffset:  0,
	}
}

func (b *InMemBlockStorage) blockStorage() {}

func (b *InMemBlockStorage) growBlocks() {
	if b.blockLen > len(b.blocks) {
		for i := len(b.blocks); i < b.blockLen; i++ {
			b.blocks = append(b.blocks, newInMemBlock(b.blockSize))
		}
	}
}

func (b *InMemBlockStorage) Write(p []byte) (n int, err error) {
	if len(p) > b.blockSize {
		return 0, SizeExceeded
	}
	remaining := b.seqWriteOffset - b.blockSize
	writeSize := len(p)

	if remaining < writeSize {
		b.seqWriteOffset = 0
		b.seqWriteBlock++
	}

	if b.blockLen < b.seqWriteBlock {
		b.blockLen = b.seqWriteBlock
		b.growBlocks()
	}

	//copy(b, p)
	nCopied := copy(b.blocks[b.seqWriteBlock].data[b.seqWriteOffset:], p)

	return nCopied, nil
}

func (b *InMemBlockStorage) Read(p []byte) (n int, err error) {
	return
}

func (b *InMemBlockStorage) Flush() error {
	return nil
}

func (b *InMemBlockStorage) Close() error {
	return nil
}

func (b *InMemBlockStorage) Sync() error {
	return nil
}

func (b *InMemBlockStorage) NumBlocks() uint {
	return 0
}

func (b *InMemBlockStorage) ReadBlock(index uint, buffer *bytes.Buffer) (n int, err error) {
	return
}

func (b *InMemBlockStorage) WriteBlock(index uint, buffer *bytes.Buffer) (n int, err error) {
	return
}

func (b *InMemBlockStorage) BlockSize() uint {
	return b.blockSize
}
