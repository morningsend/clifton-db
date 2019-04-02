package kvserver

import (
	"context"
	"github.com/zl14917/MastersProject/api/kv-client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcKVService struct {
	logger      *zap.Logger
	requests    chan bool
	concurrency int
}

func NewGrpcKVService(numConcurrentReq int) kv_client.KVStoreServer {
	if numConcurrentReq < 1 {
		numConcurrentReq = 1
	}

	srv := &GrpcKVService{
		logger:      zap.NewExample(),
		requests:    make(chan bool, numConcurrentReq),
		concurrency: numConcurrentReq,
	}

	return srv
}

func (s *GrpcKVService) Register(grpcServer *grpc.Server) {
	kv_client.RegisterKVStoreServer(grpcServer, s)
}

func (s GrpcKVService) Get(ctx context.Context, get *kv_client.GetReq) (*kv_client.Value, error) {
	s.requests <- true
	defer func() { <-s.requests }()

	s.logger.Info("GET", zap.String("key", get.Key))
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return &kv_client.Value{Key: get.Key, Value: []byte("value is here")}, nil
	}
}

func (s GrpcKVService) Put(ctx context.Context, put *kv_client.PutReq) (*kv_client.PutRes, error) {
	s.requests <- true
	defer func() { <-s.requests }()

	s.logger.Info("PUT", zap.String("key", put.Key), zap.String("value", string(put.Value)))
	return &kv_client.PutRes{Success: true}, nil
}

func (s *GrpcKVService) Delete(ctx context.Context, req *kv_client.DelReq) (*kv_client.DelRes, error) {
	s.requests <- true
	defer func() { <-s.requests }()

	s.logger.Info("DELETE", zap.String("key", req.Key))
	return &kv_client.DelRes{Ok: true}, nil
}
