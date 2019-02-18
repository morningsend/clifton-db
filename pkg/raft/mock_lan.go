package raft

import (
	"github.com/zl14917/MastersProject/pkg/raft/rpc"
	"sync"
)

type LAN struct {
	connections map[ID]map[ID]chan rpc.Message
	done chan bool
	shutdownOnce sync.Once
}

func CreateFullyConnected(peerIds []ID, bufferSize int) *LAN {

	network := make(map[ID]map[ID]chan rpc.Message)

	for _, id1 := range peerIds {
		for _, id2 := range peerIds {

			if id1 == id2 {
				continue
			}

			neighbours, ok := network[id1]
			if !ok {
				neighbours = make(map[ID]chan rpc.Message)
				network[id1] = neighbours
			}

			_, ok = neighbours[id2]
			if !ok {
				neighbours[id2] = make(chan rpc.Message)
			}
		}
	}

	lan := &LAN{
		done:        make(chan bool),
		connections: network,
	}

	return lan
}

func (lan *LAN) GetConnection(n1, n2 ID) (conn chan rpc.Message, ok bool) {
	neighbours, ok := lan.connections[n1]
	if !ok {
		return
	}

	conn, ok = neighbours[n2]
	return
}

func (lan *LAN) GetMulticastConns(node ID) (conns map[ID]chan rpc.Message, ok bool) {
	conns, ok = lan.connections[node]
	return
}

func (lan *LAN) Close() {
	go lan.shutdownOnce.Do(func(){
		for _, conns := range lan.connections {
			for _, conn := range conns {
				close(conn)
			}
		}
	})
}