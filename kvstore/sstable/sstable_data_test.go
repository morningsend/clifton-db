package sstable

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"testing"
)

func TestSstableDataWriter(t *testing.T) {
	const (
		blockSize = 4096
	)
	var (
		storage = blockstore.NewInMemBlockStorage(blockSize)
	)
	dataWriterTest(t, storage)

}

func dataWriterTest(t *testing.T, storage blockstore.BlockStorage) {
	var (
		err        error
		pos        blockstore.Position
		dataWriter = newSStableDataWriter(storage)
		_          = newSSTableDataReader(storage)
		value      = "{\"message\": \"hello_world\"}"
	)

	valueLarge := make([]byte, dataWriter.MaxValueSize+1, dataWriter.MaxValueSize+1)
	valueFit := make([]byte, dataWriter.MaxValueSize, dataWriter.MaxValueSize)

	pos, err = dataWriter.WriteValue(valueLarge)
	if err == nil {
		t.Errorf("writing a value too large should return error")
	}

	pos, err = dataWriter.WriteValue(valueFit)
	if err != nil {
		t.Error(err)
		return
	}

	if pos.Block != 1 || pos.Offset != 0 {
		t.Errorf("position is wrong. got %v", pos)
	}

	pos, err = dataWriter.WriteValue([]byte(value))
	if err != nil {
		t.Error(err)
		return
	}

	if pos.Block != 2 || pos.Offset != 0 {
		t.Errorf("position is wrong. got %v", pos)
	}

	err = dataWriter.Commit()
	if err != nil {
		t.Error(err)
		return
	}

	if dataWriter.Storage.NumBlocks() != 3 {
		t.Errorf("num blocks should be 3 but got %d", dataWriter.Storage.NumBlocks())
	}
}
