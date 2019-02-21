package raft

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"
)

func createLAN(peerIds []ID) *LAN {
	return CreateFullyConnected(peerIds, 0)
}

func createRaftCluster(peerIds []ID, lan *LAN) ([]*Cluster, []RaftFSM) {
	clusterSize := len(peerIds)
	clusters := make([]*Cluster, clusterSize, clusterSize)
	fsms := make([]RaftFSM, clusterSize, clusterSize)
	for i := 0; i < clusterSize; i++ {
		clusters[i] = NewCluster(peerIds[i], peerIds)
	}

	for i := 0; i < clusterSize; i++ {
		nodeId := peerIds[i]
		channels, ok := lan.GetMulticastConns(nodeId)
		if !ok {
			os.Exit(1)
		}
		comms := NewChannelComms(nodeId, channels)
		comms.Start()
		fsms[i] = NewRaftFSM(peerIds[i], peerIds, 10+rand.Int()%5, 20+rand.Int()%5, comms)
	}

	return clusters, fsms
}

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

	i := 34
	err := fsm.ReplicateToLog(i)
	if err != nil {
		t.Error("should be replicated")
	}
	fsm.Tick()
	term, index := fsm.Log().GetLastLogTermIndex()

	if term != 1 && index != 1 {
		t.Error("first log entry should have term 1", "index 1")
	}

	entry, ok := fsm.Log().Read(1)
	if !ok {
		t.Error("log should exist")
	}

	if entry.Data != i {
		t.Error("replicated data should be", i)
	}
}

func TestClusterOf3(t *testing.T) {
	const N = 3
	peers := []ID{1, 2, 3}
	fsms := make([]RaftFSM, N, N)
	drivers := make([]RaftDriver, N, N)
	lan := createLAN(peers)
	for i, peer := range peers {
		channels, ok := lan.GetMulticastConns(peer)
		if !ok {
			t.Error("channel should exist")
			return
		}

		comm := NewChannelComms(peer, channels)
		comm.Start()
		fsms[i] = NewRaftFSM(peer, peers, 3+2*i, 8+2*i, comm)
		drivers[i] = NewTestDriver()
		drivers[i].Init(fsms[i], time.Millisecond*5)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	for _, d := range drivers {
		go d.Run(ctx)
	}

	log := 100

	time.Sleep(time.Millisecond * 250)
	for _, raft := range fsms {
		if raft.Role() == Leader {
			_ = raft.ReplicateToLog(log)
		}
	}

	<-ctx.Done()
	lan.Close()
	for _, fsm := range fsms {
		entry, ok := fsm.Log().Read(1)
		if !ok {
			t.Error("log should exist at index", 1)
		}
		if entry.Data.(int) != log {
			t.Error("entry should be", log)
		}
	}
}
