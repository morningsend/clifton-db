package server

import (
	"github.com/morningsend/clifton-db/pkg/storage"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer  *grpc.Server
	raftStorage *storage.RaftStorage
	raftServer  *storage.RaftStorageServer
}

func New(options Options) *Server {
	return &Server{}
}

func (s *Server) Run() error {
	return nil
}
