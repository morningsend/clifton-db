package cluster_manager

import (
	"go.etcd.io/etcd/raft/raftpb"
	"sync"
)

type ClusterManagerDataStore struct {
	lock sync.RWMutex
}

func NewClusterManagerDataStore() *ClusterManagerDataStore {
	return nil
}

func (c *ClusterManagerDataStore) Snapshot() ([]byte, error) {
	return nil, nil
}

func (c *ClusterManagerDataStore) ApplyConfChange(change raftpb.ConfChange) error {
	return nil
}
