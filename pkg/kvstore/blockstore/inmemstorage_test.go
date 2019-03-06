package blockstore

import "testing"

func TestInMemBlockStorage_Allocate(t *testing.T) {

}

func TestInMemBlockStorage_SequentialWrites(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	seqWriteTest(s, t)
}

func TestInMemBlockStorage_SequentialReads(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	seqReadTest(s, t)
}

func TestInMemBlockStorage_RandomReads(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	randomReadTest(s, t)
}

func TestInMemBlockStorage_RandomWrites(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	randomWriteTest(s, t)
}
