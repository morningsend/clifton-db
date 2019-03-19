package sstable

import (
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
		position = blockstore.Position{}
	)

	indexWriter := newIndexWriter(storage)

	err := indexWriter.WriteIndex(key, false, position)
	if err != nil {
		t.Errorf("error writing index key: %s", string(key))
	}
}
