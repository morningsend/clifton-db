package raft

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/raft/rpc"
)

type Comms interface {
	BroadcastRpc(ctx context.Context, msg rpc.Message)
	Rpc(ctx context.Context, id ID, msg rpc.Message)
	Reply() <-chan rpc.Message
}

func ConnectToCluster(cluster *Cluster, virtualLan* LAN) (Comms, error) {
	conns, ok := virtualLan.GetMulticastConns(cluster.SelfID)
	if !ok {
		return nil, fmt.Errorf("can't connect to cluster")
	}

	return NewChannelComms(conns), nil
}

type TCPNetworkComms struct {

}

type ChannelComms struct {
	broadcastChannels chan interface{}
	rpcChannels map[ID] chan interface{}
	replyChannel chan interface{}
}

func NewChannelComms(conns map[ID] chan interface{}) *ChannelComms {
	return &ChannelComms {
		broadcastChannels:nil,
		rpcChannels: conns,
	}
}

func (comms *ChannelComms) BroadcastRpc(ctx context.Context, msg rpc.Message) {

}

func (comms *ChannelComms) Rpc(ctx context.Context, id ID, msg rpc.Message) {

}

func (comms *ChannelComms) Reply() <- chan rpc.Message{
	return nil
}