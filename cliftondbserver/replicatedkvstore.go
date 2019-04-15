package cliftondbserver

import (
	"bytes"
	"context"
	"errors"
	"github.com/zl14917/MastersProject/api/internal_request"
	"github.com/zl14917/MastersProject/api/kv-client"
	"github.com/zl14917/MastersProject/kvstore"
	"github.com/zl14917/MastersProject/kvstore/types"
	"go.etcd.io/etcd/raft/raftpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
	"time"
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

type ReplicatedKvStore struct {
	r *RaftNode

	proposeC    chan<- string
	commitC     <-chan *string
	confChangeC chan<- raftpb.ConfChange
	kvStore     kvstore.CliftonDBKVStore

	kvGrpcApiServer kv_client.KVStoreServer
	grpcServer      *grpc.Server

	listener *StoppableListener

	logger *zap.Logger

	reqBuilder *RequestBuilder

	requestsInFlight sync.Map

	requestTimeout time.Duration
}

func (s *ReplicatedKvStore) ServeKvStoreApi() {
	s.kvGrpcApiServer = NewGrpcKVService(100)
	kv_client.RegisterKVStoreServer(s.grpcServer, s.kvGrpcApiServer)
}

func (p *ReplicatedKvStore) NewReplicatedKvStore(conf ClusterConfig) (*ReplicatedKvStore, error) {
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

		requestTimeout: time.Second * 5000,
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

func (kv *ReplicatedKvStore) RaftHandler() http.Handler {
	return kv.r.transport.Handler()
}

func (kv *ReplicatedKvStore) KvGrpcServer() kv_client.KVStoreServer {
	return kv.kvGrpcApiServer
}

func (kv *ReplicatedKvStore) handleRequest(ctx context.Context, req *internal_request.InternalRequest) (
	*internal_request.InternalResponse, error) {

	data, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, kv.requestTimeout)

	defer cancel()

	err = kv.r.Propose(ctx, data)
	if err != nil {
		return nil, err
	}

	select {
	case <-reqCtx.Done():
		return nil, reqCtx.Err()
	default:

	}
	return nil, nil
}

func (kv *ReplicatedKvStore) ProposePut(ctx context.Context, key []byte, value []byte) (
	*internal_request.InternalResponse, error) {
	putReq := kv.reqBuilder.NewPutRequest(key, value)
	return kv.handleRequest(ctx, putReq)
}

func (st *ReplicatedKvStore) ProposeDelete(ctx context.Context, key []byte) (
	*internal_request.InternalResponse, error) {
	deleteReq := st.reqBuilder.NewDeleteRequest(key)
	return st.handleRequest(ctx, deleteReq)
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
