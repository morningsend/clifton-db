package main

import (
	"fmt"
	"github.com/zl14917/MastersProject/pkg/raft"

)

const Nodes = 3

func createLAN(peerIds []raft.ID) *raft.LAN {
	return raft.CreateFullyConnected(peerIds, 0)
}

func createRaftCluster(peerIds []raft.ID, lan *raft.LAN) ([]*raft.Cluster, []raft.RaftFSM) {
	clusterSize := len(peerIds)
	clusters := make([]*raft.Cluster, clusterSize, clusterSize)

	for i:= 0; i < clusterSize; i++ {
		clusters[i] = raft.NewCluster(peerIds[i], peerIds)
	}

	return clusters, nil
}
func main() {

	peerIds := []raft.ID{1, 2, 3}
	virtualLAN := createLAN(peerIds)
	_, _ = createRaftCluster(peerIds, virtualLAN)

	fmt.Println("starting simulation.")

	defer virtualLAN.Close()

	fmt.Println("shut down machine")
}
