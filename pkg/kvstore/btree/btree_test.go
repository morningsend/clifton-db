package btree

import (
	"fmt"
	"testing"
)

func TestBTreePut(t *testing.T) {
	tree := NewBTreeMap(2, 128)
	tree.Put("1", []byte("adsf"))
	value, ok := tree.Get("1")
	if !ok {
		t.Error("'1' should exist")
	}

	fmt.Println(string(value))
}
