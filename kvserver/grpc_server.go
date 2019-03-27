package kvserver

import (
	"context"
	"github.com/zl14917/MastersProject/api/kv-client"
	"google.golang.org/grpc"
	"log"
	"os"
)

type GrpcKVServer struct {
	logger *log.Logger
}

func NewGrpcKVServer() *GrpcKVServer {
	return &GrpcKVServer{
		logger: log.New(os.Stdout, "[cliftondb-rpc]", log.LstdFlags),
	}
}

func (s *GrpcKVServer) Register(grpcServer *grpc.Server) {
	kv_client.RegisterKVStoreServer(grpcServer, *s)
}

func (s GrpcKVServer) NewSession(context.Context, *kv_client.Client) (*kv_client.Session, error) {
	return &kv_client.Session{SessionId: "1234", Nodes: []*kv_client.Address{}}, nil
}

func (s GrpcKVServer) Get(ctx context.Context, get *kv_client.GetReq) (*kv_client.Value, error) {
	s.logger.Printf("GET %s\n", get.Key)
	return &kv_client.Value{Key: get.Key, Value: []byte("value is here")}, nil
}

func (s GrpcKVServer) Put(ctx context.Context, put *kv_client.PutReq) (*kv_client.PutRes, error) {
	s.logger.Printf("PUT %s:%s\n", put.Key, put.Value)
	return &kv_client.PutRes{Success: true}, nil
}
