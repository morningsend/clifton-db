package kvserver

import (
	"context"
	"github.com/zl14917/MastersProject/kvstore"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/wal"
)

type RaftKVStoreFSM struct {
	kvstore *kvstore.CliftonDBKVStore

	raft    *raft.Node
	raftWal *wal.WAL
}

func NewKVStoreFSM() *RaftKVStoreFSM {
	return &RaftKVStoreFSM{
	}
}

func (m *RaftKVStoreFSM) ReceiveMessage(ctx context.Context, msg interface{}) {

}

func (m *RaftKVStoreFSM) Loop() {

}
