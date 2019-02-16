package raft

import "testing"

func TestSingleRaftFsm(t *testing.T) {
	fsm := NewRaftFSM(1, 2,4)
	if fsm.Role() != Follower {
		t.Error("initial role must be follower")
	}
}