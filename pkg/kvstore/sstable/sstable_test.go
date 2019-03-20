package sstable

import (
	"bytes"
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"io/ioutil"
	"os"
	"testing"
)

const tmpDir = "/tmp/"

type testWithFileIO func(t *testing.T, dirPath string)

func WithTempDir(t *testing.T, io testWithFileIO) {

	name, err := ioutil.TempDir(tmpDir, "cliftondbtests")

	if err != nil {
		t.Fatal("can't create test directory", err)
	}

	io(t, name)

	_ = os.Remove(name)
}

func TestNewSSTable(t *testing.T) {
	WithTempDir(t, func(t *testing.T, dirPath string) {
		var options = defaultSSTableOpenOptions

		options.InMemStore = true

		sstable := NewSSTable(dirPath, &options)

		writer, err := sstable.NewWriter()
		if err != nil {
			t.Error("", err)
			return
		}

		err = writer.Write([]byte("hello"), []byte("world"), false)
		if err != nil {
			t.Error("error writing key:'hello', value: 'world'")
			return
		}

	})
}

func TestSSTableBlockIndexWriter(t *testing.T) {
	const blocksize = 4096
	var (
		storage  = blockstore.NewInMemBlockStorage(blocksize)
		key      = []byte("greetings.message")
		key2     = []byte("gzzzzz.message")
		position = blockstore.Position{Block: 12, Offset: 12}
	)

	indexWriter := newSSTableIndexWriter(storage)

	err := indexWriter.WriteIndex(key, false, position)
	if err != nil {
		t.Errorf("error writing index key: %s", string(key))
	}
	err = indexWriter.WriteIndex(key2, true, position)

	if err != nil {
		t.Error(err)
	}

	err = indexWriter.Commit()

	if err != nil {
		t.Errorf("error committing index writer: %v", err)
	}

	if storage.NumBlocks() != 2 {
		t.Errorf("two blocks at minimum: got %d", storage.NumBlocks())
	}

	indexReader := newSSTableIndexReader(storage)
	err = indexReader.ReadHeader()

	if err != nil {
		t.Error(err)
	}
	entry, ok, err := indexReader.FindIndexForKey(key)

	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Errorf("key should exist after being written")
	}

	if entry == nil {
		return
	}
	var entryPosition blockstore.Position

	entryPosition.DecodeUint64(entry.DataFileOffSet)

	if entryPosition != position {
		t.Error("error, data file offset is not same.")
	}

	if bytes.Compare(key, entry.LargeKey) != 0 {
		t.Errorf("key not the same, written %s, read %s", string(key), string(entry.LargeKey))
	}

	entry, ok, err = indexReader.FindIndexForKey(key2)

	if entry == nil {
		t.Error("key 2 should exist:", key2)
	}

}
