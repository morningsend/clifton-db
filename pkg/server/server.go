package server

import (
	"os"

	"github.com/morningsend/clifton-db/pkg/storage"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer  *grpc.Server
	raftStorage *storage.RaftStorage
	raftServer  *storage.RaftStorageServer

	sigs chan os.Signal
}

func New(options Options) *Server {
	return &Server{
		sigs: make(chan os.Signal, 10),
	}
}

func (s *Server) StartUp() error {
	return nil
}

func (s *Server) Run() error {
	return nil
}
