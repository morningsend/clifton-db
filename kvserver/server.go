package kvserver

import (
	"context"
	"github.com/zl14917/MastersProject/kv-client-api"
	"google.golang.org/grpc"
	"log"
	"os"
)

type KVServer struct {
	logger *log.Logger
}

func NewServer() *KVServer {
	return &KVServer{
		logger: log.New(os.Stdout, "[cliftondb]", log.LstdFlags),
	}
}

func (s *KVServer) Register(grpcServer *grpc.Server) {
	kv_client_api.RegisterKVStoreServer(grpcServer, *s)
}

func (s KVServer) NewSession(context.Context, *kv_client_api.Client) (*kv_client_api.Session, error) {
	return &kv_client_api.Session{SessionId: "1234", Nodes: []*kv_client_api.Address{}}, nil
}

func (s KVServer) Get(ctx context.Context, get *kv_client_api.GetReq) (*kv_client_api.Value, error) {
	s.logger.Printf("GET %s\n", get.Key)
	return &kv_client_api.Value{Key: get.Key, Value: []byte("value is here")}, nil
}

func (s KVServer) Put(ctx context.Context, put *kv_client_api.PutReq) (*kv_client_api.PutRes, error) {
	s.logger.Printf("PUT %s:%s\n", put.Key, put.Value)
	return &kv_client_api.PutRes{Success: true}, nil
}
