package cliftondbserver

import (
	"context"
	"go.etcd.io/etcd/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/etcdserver/api/snap"
	"go.etcd.io/etcd/etcdserver/api/v2stats"
	"go.etcd.io/etcd/pkg/types"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/wal"
	"go.etcd.io/etcd/wal/walpb"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
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
	defaultSnapshotLag  int      = 10000
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
	readStateC  chan<- raft.ReadState

	// channel to report commit and error
	commitC chan *string
	errorC  chan error

	Id types.ID

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

	Node                  raft.Node
	raftMemStorage        *raft.MemoryStorage
	raftPersistentStorage *raftPersistentStorage

	snapshotLag int
	snapCount   uint64
	transport   *rafthttp.Transport

	stopc     chan struct{}
	httpdonec chan struct{}
	httpstopc chan struct{}

	logger *zap.Logger

	RaftURL    string
	raftServer *http.Server

	walDirPath      string
	snapshotDirPath string
}

type ClusterConfig struct {
	SelfId      uint
	JoinCluster bool
	Peers       []PeerEntry
	StoragePath string
}

func RaftStandaloneConfig(id uint, path string) ClusterConfig {
	return ClusterConfig{
		SelfId:      id,
		JoinCluster: false,
		StoragePath: path,
	}
}

func RaftClusterConfig(id uint, peers []PeerEntry, path string) ClusterConfig {
	return ClusterConfig{
		SelfId:      id,
		JoinCluster: true,
		Peers:       peers,
		StoragePath: path,
	}
}

func (r *RaftNode) ErrorC() chan error {
	return r.errorC
}

func (r *RaftNode) CommitC() chan *string {
	return r.commitC
}

func NewRaftNode(conf ClusterConfig, proposeC <-chan string, confChangeC <-chan raftpb.ConfChange,
	options ...Options) (*RaftNode, error) {
	var err error
	lg := zap.NewExample()
	commitC := make(chan *string)
	errorC := make(chan error)

	walDir := path.Join(conf.StoragePath, walPath)
	snapDir := path.Join(conf.StoragePath, snapshotPath)

	n := &RaftNode{
		ClusterId:      defaultClusterID,
		Id:             types.ID(conf.SelfId),
		Peers:          []PeerEntry{},
		tickInterval:   defaultTickInterval,
		commitC:        commitC,
		errorC:         errorC,
		stopc:          make(chan struct{}),
		proposeC:       proposeC,
		confChangeC:    confChangeC,
		logger:         lg,
		raftMemStorage: raft.NewMemoryStorage(),

		snapshotLag: defaultSnapshotLag,

		walDirPath:      walDir,
		snapshotDirPath: snapDir,
	}

	for _, option := range options {
		option.apply(n)
	}
	var (
		WAL         *wal.WAL
		hardState   raftpb.HardState
		entries     []raftpb.Entry
		walsnapshot walpb.Snapshot
		snapshot    *raftpb.Snapshot
	)

	snapshotter := snap.New(zap.NewExample(), snapDir)
	snapshot, err = snapshotter.Load()

	if err != nil {
		lg.Fatal("error loading snapshot",
			zap.String("snapshotDirPath", snapDir),
			zap.Error(err),
		)
	}

	walsnapshot.Term, walsnapshot.Index = snapshot.Metadata.Term, snapshot.Metadata.Index

	if !wal.Exist(walDir) {
		WAL, err = wal.Create(lg, walDir, []byte{})
		if err != nil {
			lg.Fatal("failed to create wal",
				zap.String("walDirPath", walDir),
				zap.Error(err),
			)
		}
		if err = WAL.Close(); err != nil {
			lg.Fatal("error closing WAL", zap.Error(err))
		}
	}

	WAL, hardState, entries = readWALContent(lg, walDir, walsnapshot)

	n.raftPersistentStorage = NewRaftPersistentStorage(WAL, snapshotter)

	if err = n.raftMemStorage.ApplySnapshot(*snapshot); err != nil {
		lg.Fatal("error applying snapshot", zap.Error(err))
	}

	if err = n.raftMemStorage.SetHardState(hardState); err != nil {
		lg.Fatal("error setting raft hardstate", zap.Error(err))
	}

	if err = n.raftMemStorage.Append(entries); err != nil {
		lg.Fatal("error appending raft entries to storage", zap.Error(err))
	}

	n.transport = &rafthttp.Transport{
		Logger:      zap.NewExample(),
		ID:          types.ID(n.Id),
		ClusterID:   n.ClusterId,
		Raft:        n,
		ServerStats: v2stats.NewServerStats("", ""),
		LeaderStats: v2stats.NewLeaderStats(strconv.Itoa(int(n.Id))),
	}

	return n, nil
}

func (r *RaftNode) ProcessEntries(entry []raftpb.Entry) bool {
	return true
}

func (r *RaftNode) Start() {
	var err error
	peers := make([]raft.Peer, 0)
	c := &raft.Config{
		ID:              uint64(r.Id),
		ElectionTick:    10,
		HeartbeatTick:   2,
		Storage:         r.raftMemStorage,
		MaxInflightMsgs: 128,
		MaxSizePerMsg:   1024 * 1024,
		CheckQuorum:     true,
		PreVote:         true,
	}
	r.Node = raft.StartNode(c, peers)
	err = r.transport.Start()
	if err != nil {
		r.logger.Fatal("error starting raft transport", zap.Error(err))
	}
}

func (r *RaftNode) Restart() {

}
func (r *RaftNode) Loop(ticker *time.Ticker) error {
	var (
		snapshot, err = r.raftMemStorage.Snapshot()
	)

	r.confState = snapshot.Metadata.ConfState
	r.snapshotIndex = snapshot.Metadata.Index
	r.appliedIndex = snapshot.Metadata.Index

	if err != nil {
		r.logger.Error("error reading snapshot from storage", zap.Error(err))
		return err
	}
	tickerChan := ticker.C
	var isLeader = false
	for {
		select {
		case <-tickerChan:
			r.Node.Tick()
		case rd := <-r.Node.Ready():

			if len(rd.ReadStates) != 0 {
				select {
				case r.readStateC <- rd.ReadStates[len(rd.ReadStates)-1]:
				case <-r.stopc:
					return nil
				}
			}

			isLeader = rd.RaftState == raft.StateLeader

			if isLeader {
				r.transport.Send(rd.Messages)
			}

			if err := r.raftPersistentStorage.Save(rd.HardState, rd.Entries); err != nil {
				if r.logger != nil {
					r.logger.Fatal("error saving entries to wal", zap.Error(err))
				}
			}

			if !raft.IsEmptySnap(rd.Snapshot) {

				if err = r.raftPersistentStorage.SaveSnap(rd.Snapshot); err != nil {
					r.logger.Error("error saving snapshot", zap.Error(err))
				}
				if err = r.raftMemStorage.ApplySnapshot(rd.Snapshot); err != nil {
					r.logger.Error("error applying snapshot", zap.Error(err))
				}

				if err = r.PublishSnapshot(rd.Snapshot); err != nil {
					r.logger.Error("error publishing snapshot", zap.Error(err))
				}
			}

			err = r.raftMemStorage.Append(rd.Entries)

			if !isLeader {
				msg := rd.Messages
				r.transport.Send(msg)
			}

			r.ProcessEntries(rd.CommittedEntries)
			r.Node.Advance()

		case err := <-r.transport.ErrorC:
			r.logger.Error("transport: error sending message", zap.Error(err))
		case <-r.stopc:
			tickerChan = nil
			r.Stop()
			return nil
		}
	}
}

func (r *RaftNode) applyEntries(entries []raftpb.Entry) {

}

func (r *RaftNode) TriggerSnapshot() {
	var err error
	//get data from kv store.
	snapshot, err := r.raftMemStorage.CreateSnapshot(r.appliedIndex, &r.confState, nil)
	if err != nil {
		if err == raft.ErrSnapOutOfDate {
			return
		}
		r.logger.Fatal("cannot create snapshot", zap.Error(err))
	}

	err = r.raftPersistentStorage.SaveSnap(snapshot)
	if err != nil {
		r.logger.Fatal("error saving snapshot to persistent storage", zap.Error(err))
	}
	
}

func (r *RaftNode) Propose(ctx context.Context, entryData []byte) error {
	return r.Node.Propose(ctx, entryData)
}

func (r *RaftNode) ProposeConfChange(ctx context.Context, change raftpb.ConfChange) error {
	return r.Node.ProposeConfChange(ctx, change)
}
func (r *RaftNode) Stop() {
	close(r.commitC)
	close(r.errorC)
	r.transport.Stop()
	r.Node.Stop()
}

func (r *RaftNode) PublishSnapshot(snapshot raftpb.Snapshot) error {
	return nil
}

func (r *RaftNode) PublishEntries(entries []raftpb.Entry) bool {
	return true
}

func (r *RaftNode) StartRaftServer() error {
	raftUrl, err := url.Parse(r.RaftURL)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", raftUrl.Host)
	if err != nil {
		return err
	}

	sl, err := NewStoppableListener(listener)
	if err != nil {
		r.logger.Error("error creating stoppable listener", zap.Error(err))
		cerr := listener.Close()
		if cerr != nil {
			r.logger.Error("error closing listener", zap.Error(cerr))
		}
		return err
	}

	r.raftServer = &http.Server{
		Handler: r.transport.Handler(),
	}

	err = r.raftServer.Serve(sl)
	return err
}
func (rc *RaftNode) Process(ctx context.Context, m raftpb.Message) error {
	return rc.Node.Step(ctx, m)
}
func (rc *RaftNode) IsIDRemoved(id uint64) bool                           { return false }
func (rc *RaftNode) ReportUnreachable(id uint64)                          {}
func (rc *RaftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}
