package blockstore

import (
	"bytes"
	"io"
	"testing"
)

func fillArray(arr []byte, value byte) {
	for i, _ := range arr {
		arr[i] = value
	}
}

func seqWriteTest(s BlockStorage, t *testing.T) {
	size := s.BlockSize()

	buf := make([]byte, size-1, size-1)
	bufTooLarge := make([]byte, size+1, size+1)

	n, err := s.Write(buf)
	if n != len(buf) {
		t.Error("number of bytes written is less than ")
	}

	if err != nil {
		t.Error("error writing data to block storage", err)
	}

	if s.NumBlocks() != 1 {
		t.Error("after first write, number of blocks should be 1")
	}
	_, err = s.Write(bufTooLarge)

	if err == nil {
		t.Errorf("writing a block too large should fail, but succeeded. block size %d, buffer size %d",
			size,
			len(bufTooLarge),
		)
	}

	_, _ = s.Write(buf)
	if s.NumBlocks() != 2 {
		t.Error("sequential write should grow block.")
	}
}

func seqReadTest(s BlockStorage, t *testing.T) {
	const value1 = 0xf1
	const value2 = 0xf2
	size := s.BlockSize()
	buf1 := make([]byte, size, size)
	buf2 := make([]byte, size, size)
	fillArray(buf1, value1)
	fillArray(buf2, value2)

	s.Write(buf1)
	s.Write(buf2)

	if s.NumBlocks() != 2 {
		t.Error("should have 2 blocks")
		return
	}

	readBuffer := make([]byte, size/2, size/2)
	nRead, err := s.Read(readBuffer)

	if nRead != len(readBuffer) {
		t.Error("number of bytes read is different than size of buffer", size/2, nRead)
		return
	}

	if err != nil {
		t.Error("error reading from block storage", err)
		return
	}

	if !bytes.Equal(readBuffer, buf1[:len(readBuffer)]) {
		t.Error("bytes read should equal to bytes written")
	}

	_, _ = s.Read(readBuffer)
	_, _ = s.Read(readBuffer)
	_, _ = s.Read(readBuffer)

	_, err = s.Read(readBuffer)

	if err != io.EOF {
		t.Error("end of block should return EOF error, but got", err)
	}

}

func randomReadTest(s BlockStorage, t *testing.T) {

}

func randomWriteTest(s BlockStorage, t *testing.T) {
	const (
		Blocks = 2
		Value1 = 0xf1
		Value2 = 0xf2
	)
	var BlockSize = s.BlockSize()

	nblock, err := s.Allocate(Blocks)
	if nblock != Blocks {
		t.Errorf("should allocate block %d but got %d", Blocks, nblock)
	}

	if err != nil {
		t.Error(err)
		return
	}

	data1 := make([]byte, BlockSize, BlockSize)
	data2 := make([]byte, BlockSize, BlockSize)

	fillArray(data1, Value1)
	fillArray(data2, Value1)
	buf1 := bytes.NewBuffer(data1)
	buf2 := bytes.NewBuffer(data2)
	n, err := s.WriteBlock(0, buf1)

	if err != nil {
		t.Error(err)
		return
	}

	n, err := s.WriteBlock(1, buf2)
}
