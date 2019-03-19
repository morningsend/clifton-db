package sstable

import (
	"bytes"
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
}
