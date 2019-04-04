package cliftondbserver

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/zl14917/MastersProject/api/kv-client"
	"github.com/zl14917/MastersProject/kvstore"
	"github.com/zl14917/MastersProject/kvstore/types"
	"go.etcd.io/etcd/raft/raftpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type ReadConsistencyLevel int

const (
	Serializable ReadConsistencyLevel = iota
	Quorum
	Linearizable
)

var ConsistencyLevelNotSupportedErr = errors.New("consistency level not supported")

const (
	DefaultReadConsistency = Quorum
)

type GetOptions struct {
	ReadConsistencyLevel ReadConsistencyLevel
}

var DefaultGetOption = GetOptions{
	DefaultReadConsistency,
}

type KvStoreAction struct {
	Type    KvStoreActionType
	KeyHash uint32
	Key     string
	Value   []byte
}

type KvStoreActionType int

const (
	GetAction KvStoreActionType = iota
	PutAction
	DeleteAction
)

type ReplicatedKvStore struct {
	r           *RaftNode
	proposeC    chan<- string
	commitC     <-chan *string
	confChangeC chan<- raftpb.ConfChange
	kvStore     kvstore.CliftonDBKVStore

	kvGrpcApiServer kv_client.KVStoreServer
	grpcServer      *grpc.Server

	listener *StoppableListener

	logger *zap.Logger
}

func (s *ReplicatedKvStore) ServeKvStoreApi() {
	s.kvGrpcApiServer = NewGrpcKVService(100)
	kv_client.RegisterKVStoreServer(s.grpcServer, s.kvGrpcApiServer)
}

func (p *ReplicatedKvStore) NewReplicatedKvStore(conf RaftConfig) (*ReplicatedKvStore, error) {
	logger := zap.NewExample()
	proposeC := make(chan string)
	confChangeC := make(chan raftpb.ConfChange)
	commitC := make(chan *string)

	r, err := NewRaftNode(conf, proposeC, confChangeC)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		return nil, err
	}

	sl, err := NewStoppableListener(listener)

	if err != nil {
		logger.Error("error starting listener", zap.Error(err))
		err = listener.Close()
		return nil, err
	}

	store := &ReplicatedKvStore{
		proposeC:   proposeC,
		commitC:    commitC,
		grpcServer: grpc.NewServer(),
		listener:   sl,
		r:          r,
		logger:     logger,
	}

	p.r = r

	err = p.r.StartRaftServer()

	if err != nil {
		return nil, err
	}

	go func() {
		err := p.grpcServer.Serve(sl)
		if err != nil {
			store.logger.Error("error serving GRPC", zap.Error(err))
		}
	}()

	return store, nil
}

func (b *ReplicatedKvStore) ProposePut(ctx context.Context, key string, value []byte) error {
	var buf bytes.Buffer
	putReq := KvStoreAction{
		Key:   key,
		Value: value,
		Type:  PutAction,
	}
	if err := gob.NewEncoder(&buf).Encode(&putReq); err != nil {
		return err
	}

	return b.r.Propose(ctx, buf.Bytes())
}

func (st *ReplicatedKvStore) ProposeDelete(ctx context.Context, key string) error {
	var buf bytes.Buffer
	delReq := KvStoreAction{
		Key:  key,
		Type: DeleteAction,
	}
	if err := gob.NewEncoder(&buf).Encode(&delReq); err != nil {
		return err
	}

	return st.r.Propose(ctx, buf.Bytes())
}

func (st *ReplicatedKvStore) Get(ctx context.Context, key string, options GetOptions) ([]byte, error) {
	switch options.ReadConsistencyLevel {
	case Serializable:
		value, ok, err := st.kvStore.Get(types.KeyType(key))
		if err != nil || !ok {
			return nil, err
		}
		return value, nil
	case Quorum:
		return nil, nil
	case Linearizable:
		return nil, nil
	default:
		return nil, ConsistencyLevelNotSupportedErr
	}
}

func (st *ReplicatedKvStore) ProposeC() chan<- string {
	return st.proposeC
}

func (s *ReplicatedKvStore) Stop() {
	s.listener.Stop()
	s.r.Stop()
}

func (s *ReplicatedKvStore) Apply(entryData []byte) {
	_ = bytes.NewBuffer(entryData)
	return
}
