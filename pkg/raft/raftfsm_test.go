package raft

import "testing"

func TestSingleRaftFsm(t *testing.T) {
	fsm := NewRaftFSM(ID(1), []ID{1}, 4, 8, NewChannelComms(nil));
	if fsm.Role() != Follower {
		t.Error("initial role must be follower")
	}
}
