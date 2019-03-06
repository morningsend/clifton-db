package blockstore

import "testing"


func TestInMemBlockStorage_SequentialWrites(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	seqWriteTest(s, t)
	_ = s.Close()
}

func TestInMemBlockStorage_SequentialReads(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	seqReadTest(s, t)
	_ = s.Close()
}

func TestInMemBlockStorage_RandomReads(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	randomReadTest(s, t)
	_ = s.Close()
}

func TestInMemBlockStorage_RandomWrites(t *testing.T) {
	const BlockSize = 128
	s := NewInMemBlockStorage(BlockSize)
	randomWriteTest(s, t)
	_ = s.Close()
}
