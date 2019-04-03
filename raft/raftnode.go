package raft

import (
	"context"
	"errors"
	"go.etcd.io/etcd/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/etcdserver/api/snap"
	"go.etcd.io/etcd/etcdserver/api/v2stats"
	"go.etcd.io/etcd/pkg/types"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/wal"
	"go.uber.org/zap"
	"path"
	"strconv"
	"time"
)

const snapshotPath = "/snapshots"
const walPath = "/wal"

var defaultSnapshotCount uint64 = 10000

const (
	defaultClusterID    types.ID = 0x1000
	defaultTickInterval          = time.Millisecond * 100
)

type Options interface {
	apply(r *RaftNode)
}

type JoinClusterOption struct{}

func (j *JoinClusterOption) apply(r *RaftNode) { r.join = true }
func JoinCluster() Options                     { return &JoinClusterOption{} }

type ClusterOptions struct {
	id    uint
	peers []PeerEntry
}
type PeerEntry struct {
	Id          uint
	NetworkAddr string
}

func (o *ClusterOptions) apply(r *RaftNode)               { r.Peers = o.peers; r.Id = types.ID(o.id) }
func WithPeers(nodeId uint, entries ...PeerEntry) Options { return &ClusterOptions{id: nodeId, peers: entries} }

type DirPathOption struct{ raftDirPath string }

func (o *DirPathOption) apply(r *RaftNode) {
	r.SnapDir = path.Join(o.raftDirPath, snapshotPath)
	r.WalDir = path.Join(o.raftDirPath, walPath)
}

func WithDirPath(dirPath string) Options { return &DirPathOption{raftDirPath: dirPath} }

// Raft nodes are connected

type RaftNode struct {
	// events to apply to state machine
	proposeC    <-chan string
	confChangeC <-chan raftpb.ConfChange

	// channel to report commit and error
	commitC chan<- *string
	errorC  chan<- error

	Id        types.ID
	Peers     []PeerEntry
	ClusterId types.ID
	join      bool

	tickInterval time.Duration

	WalDir  string
	SnapDir string

	LastIndex uint64

	confState     raftpb.ConfState
	snapshotIndex uint64
	appliedIndex  uint64

	Node        raft.Node
	raftStorage *raft.MemoryStorage
	wal         *wal.WAL

	snapshotter      *snap.Snapshotter
	snapshotterReady chan *snap.Snapshotter

	snapCount uint64
	transport *rafthttp.Transport

	stopc     chan struct{}
	httpdonec chan struct{}
	httpstopc chan struct{}

	logger *zap.Logger
}

type RaftConfig struct {
	SelfId      uint
	JoinCluster bool
	Peers       []PeerEntry
}

func RaftStandaloneConfig() RaftConfig {
	return RaftConfig{
		SelfId:      1,
		JoinCluster: false,
	}
}

func RaftClusterConfig(id uint, peers []PeerEntry) RaftConfig {
	return RaftConfig{
		SelfId:      id,
		JoinCluster: true,
		Peers:       peers,
	}
}

func NewRaftNode(conf RaftConfig, proposeC <-chan *string, confChangeC <-chan raftpb.ConfChange,
	options ...Options) (*RaftNode, error) {
	var err error

	commitC := make(chan *string)
	errorC := make(chan error)

	n := &RaftNode{
		ClusterId:    defaultClusterID,
		Id:           types.ID(conf.SelfId),
		Peers:        []PeerEntry{},
		tickInterval: defaultTickInterval,
		commitC:      commitC,
		errorC:       errorC,
		stopc:        make(chan struct{}),
	}

	for _, option := range options {
		option.apply(n)
	}

	wl, err := n.replayWAL()
	if err != nil {
		return nil, err
	}

	n.wal = wl

	n.transport = &rafthttp.Transport{
		Logger:      zap.NewExample(),
		ID:          types.ID(n.Id),
		ClusterID:   n.ClusterId,
		Raft:        n,
		ServerStats: v2stats.NewServerStats("", ""),
		LeaderStats: v2stats.NewLeaderStats(strconv.Itoa(int(n.Id))),
	}

	n.logger, err = zap.NewProduction()

	if err != nil {
		return nil, err
	}

	err = n.transport.Start()
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *RaftNode) replayWAL() (*wal.WAL, error) {
	return nil, errors.New("not implemented")
}

func (r *RaftNode) Loop2(ticker *time.Ticker) {
	var confChangeCount uint = 0
	proposeC := r.proposeC
	confChangeC := r.confChangeC

	for {
		select {
		case _, ok := <-proposeC:
			r.logger.Info("received proposed msg")
			if ! ok {
				proposeC = nil
			}
		case _, ok := <-confChangeC:
			if !ok {
				confChangeC = nil
			} else {
				confChangeCount++
			}
		}
	}
}

func (r *RaftNode) ProcessEntry(entry raftpb.Entry) {

}
func (r *RaftNode) Loop(ticker *time.Ticker) error {
	var (
		snap, err = r.raftStorage.Snapshot()
	)

	r.confState = snap.Metadata.ConfState
	r.snapshotIndex = snap.Metadata.Index
	r.appliedIndex = snap.Metadata.Index

	if err != nil {
		r.logger.Error("error reading snapshot from storage", zap.Error(err))
		return err
	}

	for {
		select {
		case <-ticker.C:
			r.Node.Tick()
		case rd := <-r.Node.Ready():
			err = r.wal.Save(rd.HardState, rd.Entries)
			r.logger.Error("error saving entries to wal", zap.Error(err))

			if !raft.IsEmptySnap(rd.Snapshot) {

				if err = r.SaveSnap(rd.Snapshot); err != nil {
					r.logger.Error("error saving snapshot", zap.Error(err))
				}
				if err = r.raftStorage.ApplySnapshot(rd.Snapshot); err != nil {
					r.logger.Error("error applying snapshot", zap.Error(err))
				}

				if err = r.PublishSnapshot(rd.Snapshot); err != nil {
					r.logger.Error("error publishing snapshot", zap.Error(err))
				}
			}

			err = r.raftStorage.Append(rd.Entries)
			r.transport.Send(rd.Messages)

			for _, entry := range rd.CommittedEntries {
				r.ProcessEntry(entry)
				if entry.Type == raftpb.EntryConfChange {
					var cc raftpb.ConfChange
					if err := cc.Unmarshal(entry.Data); err != nil {
						continue
					}
					r.Node.ApplyConfChange(cc)
				}
			}
			r.Node.Advance()
		case err := <-r.transport.ErrorC:
			r.logger.Error("transport: error sending message", zap.Error(err))
		case <-r.stopc:
			r.Stop()
			return nil
		}
	}
}

func (r *RaftNode) Stop() {
	close(r.commitC)
	close(r.errorC)
	r.transport.Stop()
	r.Node.Stop()
}

func (r *RaftNode) SaveSnap(snapshot raftpb.Snapshot) error {
	return nil
}

func (r *RaftNode) PublishSnapshot(snapshot raftpb.Snapshot) error {
	return nil
}

func (r *RaftNode) PublishEntries(entries []raftpb.Entry) bool {
	return false
}

func (rc *RaftNode) Process(ctx context.Context, m raftpb.Message) error {
	return rc.Node.Step(ctx, m)
}
func (rc *RaftNode) IsIDRemoved(id uint64) bool                           { return false }
func (rc *RaftNode) ReportUnreachable(id uint64)                          {}
func (rc *RaftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}
