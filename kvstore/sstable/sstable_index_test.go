package sstable

import (
	"bytes"
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"reflect"
	"testing"
)

func TestNextMultipleOf4Uint(t *testing.T) {
	input := []uint{0, 3, 4, 14, 128}
	expect := []uint{0, 4, 4, 16, 128}

	for i, n := range input {
		result := NextMultipleOf4Uint(n)
		if result != expect[i] {
			t.Errorf("NextMultipleOf4(%d) expect %d got %d", n, expect[i], result)
		}
	}
}

func TestSSTableIndexEntry_Marshall(t *testing.T) {
	message := "hello_world_12345"
	var unmarshalledEntry SSTableIndexEntry

	entry := SSTableIndexEntry{
		Flags:          SSTableIndexKeyInsert,
		KeyLen:         uint32(len(message)),
		DataFileOffSet: 123456789,
		LargeKey:       []byte(message),
	}

	marshalledSize := MarshalledSizeOfIndexEntry(&entry)

	buffer := bytes.NewBuffer(nil)
	n, err := entry.Marshall(buffer)
	if err != nil {
		t.Errorf("error marshalling entry: %v", err)
		return
	}

	if n != int(marshalledSize) {
		t.Errorf("marshalled entry should have size %d, but got %d", marshalledSize, n)
	}

	err = unmarshalledEntry.UnMarshall(buffer)
	if err != nil {
		t.Error("error unmarshalling entry:", err)
	}

	if !reflect.DeepEqual(unmarshalledEntry, entry) {
		t.Error("marshalled entry does not equal to original value.")
	}

	if bytes.Compare(unmarshalledEntry.LargeKey, entry.LargeKey) != 0 {
		t.Error("marshalled/unmarshalled key should be identical bytes")
	}
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
