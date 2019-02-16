package rpc

import "testing"

func TestInterfaceImplemented(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	var req AppendEntriesReq
	var reqData = NewA
}
