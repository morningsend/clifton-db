package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.etcd.io/etcd/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/pkg/types"
	"go.uber.org/zap"
)

type RaftStorageServer struct {
	server    *http.Server
	transport *rafthttp.Transport
	rs        *RaftStorage

	done chan struct{}
}

func NewRaftStorageServer(bindAddr string, bindPort int, clusterId uint64, rs *RaftStorage) (*RaftStorageServer, error) {
	raftTransport := &rafthttp.Transport{
		Logger:    zap.NewExample(),
		ID:        types.ID(rs.NodeId),
		ClusterID: types.ID(clusterId),
		Raft:      rs,
		ErrorC:    make(chan error, 1),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", bindAddr, bindPort),
		Handler:      raftTransport.Handler(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	s := &RaftStorageServer{
		transport: raftTransport,
		server:    server,
		done:      make(chan struct{}),
	}

	return s, nil
}

func (t *RaftStorageServer) Start() error {
	err := t.transport.Start()
	if err != nil {
		return err
	}

	go func() {
		for err := range t.transport.ErrorC {
			log.Printf("raft transport error: %v", err)
		}
	}()

	go t.server.ListenAndServe()
	return nil
}

func (t *RaftStorageServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := t.server.Shutdown(ctx)
	t.transport.Stop()
	return err
}

func (s *RaftStorageServer) Storage() *RaftStorage {
	return s.rs
}
