package storage

import (
	"context"
	"net"
	"os"

	"go.etcd.io/etcd/etcdserver/api/rafthttp"

	"github.com/morningsend/clifton-db/pkg/storagepb"
	"go.etcd.io/etcd/raft/raftpb"

	"go.etcd.io/etcd/etcdserver/api/snap"

	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/wal"
	"go.etcd.io/etcd/wal/walpb"
	"go.uber.org/zap"
)

type NodeState struct {
	LastAppliedIndex   uint64
	LastCommittedIndex uint64
	SnapshotIndex      uint64
	ConfState          raftpb.ConfState
}

type NodeConf struct {
	NodeId          uint64
	Join            bool
	WalDirPath      string
	SnapshotDirPath string
	Peers           []Peer
}

type Peer struct {
	ID      uint64
	Context storagepb.PeerContext
}

type RaftStorage struct {
	NodeConf
	NodeState

	hasOldWal     bool
	node          raft.Node
	memoryStorage *raft.MemoryStorage
	wal           *wal.WAL
	snapshotter   *snap.Snapshotter

	logger *zap.Logger

	done chan struct{}
	stop chan struct{}
	errC chan error
}

type RaftStorageOptions struct {
	NodeId      uint64
	WalDirPath  string
	SnapDirPath string

	Join  bool
	Peers []string
}

// RaftServiceContext external dependencies which a Raft storage node
// relies on is passed in dependency injection style.
type RaftServiceContext struct {
	Listener net.Listener
	Logger   *zap.Logger
}

var _ rafthttp.Raft = (*RaftStorage)(nil)

func (r *RaftStorage) IsIDRemoved(id uint64) bool { return false }

func (r *RaftStorage) ReportUnreachable(id uint64) {}

func (r *RaftStorage) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}

func (r *RaftStorage) Process(ctx context.Context, m raftpb.Message) error {
	return r.node.Step(ctx, m)
}

func (r *RaftStorage) Snapshotter() *snap.Snapshotter {
	return r.snapshotter
}

func NewRaftStorage(options RaftStorageOptions, context *RaftServiceContext) (*RaftStorage, error) {
	s := &RaftStorage{
		done: make(chan struct{}),
		stop: make(chan struct{}),
		errC: make(chan error, 1),
	}

	s.NodeId = options.NodeId
	s.Join = options.Join
	s.WalDirPath = options.WalDirPath
	s.SnapshotDirPath = options.SnapDirPath

	err := s.LoadRaftPersistentState()
	if err != nil {
		return nil, err
	}

	err = s.startRaftNode()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *RaftStorage) Start() {
	go r.serve()
}

func (r *RaftStorage) serve() {
	defer func() {
		close(r.errC)
	}()

	for {
		select {
		case <-r.done:
			return

		}
	}
}

func (r *RaftStorage) startRaftNode() error {
	peers := make([]raft.Peer, len(r.Peers))
	for i, p := range r.Peers {
		peers[i] = raft.Peer{
			ID:      p.ID,
			Context: nil,
		}
	}

	conf := &raft.Config{
		Applied:                   r.LastAppliedIndex,
		ID:                        r.NodeId,
		ElectionTick:              10,
		HeartbeatTick:             1,
		Storage:                   r.memoryStorage,
		MaxSizePerMsg:             1024 * 1024 * 1,
		MaxInflightMsgs:           128,
		MaxUncommittedEntriesSize: 1 << 30,
	}

	if r.hasOldWal {
		r.node = raft.RestartNode(conf)
	} else {
		if r.Join {
			r.node = raft.StartNode(conf, nil)
		} else {
			r.node = raft.StartNode(conf, peers)
		}
	}

	return nil
}

func (r *RaftStorage) Close() error {
	r.SavePersistentState()
	r.node.Stop()
	return r.wal.Close()
}

func (r *RaftStorage) SavePersistentState() error {
	snapshot, err := r.memoryStorage.CreateSnapshot(r.LastAppliedIndex, &r.ConfState, nil)
	if err != nil {
		return err
	}

	walSnapshot := walpb.Snapshot{
		Term:  snapshot.Metadata.Term,
		Index: snapshot.Metadata.Index,
	}

	if err := r.wal.SaveSnapshot(walSnapshot); err != nil {
		return err
	}

	if err := r.snapshotter.SaveSnap(snapshot); err != nil {
		return err
	}

	return nil
}

func (r *RaftStorage) LoadRaftPersistentState() error {
	raftSnapshot, err := r.openSnapshot()
	if err != nil {
		return nil
	}

	oldWal, err := r.openWALAt(raftSnapshot.Metadata.Term, raftSnapshot.Metadata.Index)
	if err != nil {
		return err
	}

	// replaying the WAL log, returns etcd raft's new state
	_, state, entries, err := r.wal.ReadAll()
	if err != nil {
		return err
	}

	r.memoryStorage = raft.NewMemoryStorage()
	if err := r.memoryStorage.ApplySnapshot(*raftSnapshot); err != nil {
		return err
	}

	if err := r.memoryStorage.SetHardState(state); err != nil {
		return err
	}

	if err := r.memoryStorage.Append(entries); err != nil {
		return err
	}

	nodeState := NodeState{
		ConfState:        raftSnapshot.Metadata.ConfState,
		LastAppliedIndex: raftSnapshot.Metadata.Index,
		SnapshotIndex:    raftSnapshot.Metadata.Index,
	}

	r.NodeState = nodeState
	r.hasOldWal = oldWal
	return nil
}

func (r *RaftStorage) openSnapshot() (*raftpb.Snapshot, error) {
	s := snap.New(zap.NewExample(), r.SnapshotDirPath)
	snapshot, err := s.Load()
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (r *RaftStorage) openWALAt(term uint64, index uint64) (oldWal bool, err error) {
	var (
		wlog *wal.WAL
	)

	oldWal = wal.Exist(r.WalDirPath)

	if !oldWal {
		if err := os.Mkdir(r.WalDirPath, 0750); err != nil {
			return false, err
		}

		wlog, err = wal.Create(zap.NewExample(), r.WalDirPath, nil)
		if err != nil {
			return false, nil
		}
		wlog.Close()
	}

	walsnap := walpb.Snapshot{
		Term:  term,
		Index: index,
	}

	wlog, err = wal.Open(zap.NewExample(), r.WalDirPath, walsnap)
	if err != nil {
		return oldWal, err
	}

	return oldWal, nil
}

func (r *RaftStorage) IsLeader() bool {
	return false
}

func (r *RaftStorage) LeaderId() uint64 {
	return r.NodeId
}
