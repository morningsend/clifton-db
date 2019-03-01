package wal

import "testing"

func TestNewWAL(t *testing.T) {
	x := NewWAL("/tmp/wal/new_test")
}