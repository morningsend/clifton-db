package cluster_manager

import (
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
