package storage

import (
	"context"

	"github.com/morningsend/clifton-db/pkg/storagepb"

	"go.etcd.io/etcd/raft/raftpb"
)

type Proposer interface {
	ProposeConfChange(ctx context.Context, confChange *raftpb.ConfChange)
	ProposeLogEntry(ctx context.Context, log *storagepb.LogEntry)
}
