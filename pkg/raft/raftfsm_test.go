package raft

import "testing"

func TestSingleRaftFsm(t *testing.T) {
	fsm := NewRaftFSM(ID(1), []ID{1}, 4, 8, NewChannelComms(1, nil))
	if fsm.Role() != Follower {
		t.Error("initial role must be follower")
	}

	for i := 0; i < 8; i++ {
		fsm.Tick()
	}

	if fsm.Role() != Leader {
		t.Error("after heartbeat time should become a Leader since only one node:", fsm.Role())
	}

	fsm.AppendEntries()
}
