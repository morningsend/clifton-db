package blockstore

import (
	"bytes"
	"fmt"
	"testing"
)

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

func TestBytesBufferSizeGrow(t *testing.T) {

	buffer := bytes.NewBuffer(nil)
	buffer.Grow(10)
	buffer.Write([]byte{4, 5, 6, 7, 8, 9, 0, 1, 2, 3})
	fmt.Println(buffer.Bytes()[0:10])
	buffer.Reset()
	fmt.Println(buffer.Bytes()[0:10])

	fmt.Println(buffer.Bytes()[0:10])
	data := buffer.Bytes()[0:10]

	fmt.Println(data)
	data[3] = 99
	fmt.Println(data)
	fmt.Println(buffer.Bytes()[0:10])

	buffer.Reset()
	buffer.Write(buffer.Bytes()[0:10])
	fmt.Println(buffer.Len())
}
