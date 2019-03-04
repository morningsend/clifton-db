package wal

import "testing"

func TestNewWAL(t *testing.T) {
	_ = NewWAL("/tmp/wal/new_test", WithAutoSync(), WithCleanUp())

}
