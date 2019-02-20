package main

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/raft"
	"math/rand"
	"os"
	"runtime"
	"time"
)

const Nodes = 3

func createLAN(peerIds []raft.ID) *raft.LAN {
	return raft.CreateFullyConnected(peerIds, 0)
}

func createRaftCluster(peerIds []raft.ID, lan *raft.LAN) ([]*raft.Cluster, []raft.RaftFSM) {
	clusterSize := len(peerIds)
	clusters := make([]*raft.Cluster, clusterSize, clusterSize)
	fsms := make([]raft.RaftFSM, clusterSize, clusterSize)
	for i := 0; i < clusterSize; i++ {
		clusters[i] = raft.NewCluster(peerIds[i], peerIds)
	}

	for i := 0; i < clusterSize; i++ {
		nodeId := peerIds[i]
		channels, ok := lan.GetMulticastConns(nodeId)
		if !ok {
			os.Exit(1)
		}
		comms := raft.NewChannelComms(channels)
		fsms[i] = raft.NewRaftFSM(peerIds[i], peerIds, 10+rand.Int()%5, 20+rand.Int()%5, comms)
	}

	return clusters, fsms
}

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {

	peerIds := []raft.ID{1, 2, 3}
	virtualLAN := createLAN(peerIds)
	_, raftFsms := createRaftCluster(peerIds, virtualLAN)

	fmt.Println("starting simulation.")

	ctx, cancel := context.WithCancel(context.Background())

	for _, fsm := range raftFsms {
		driver := raft.NewTestDriver()
		driver.Init(fsm, time.Millisecond*10)
		go driver.Run(ctx)
	}

	time.Sleep(10 * time.Second)
	cancel()

	defer virtualLAN.Close()

	fmt.Println("shut down machine")
}
