package raft

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/raft/rpc"
	"sync"
	"sync/atomic"
)

type Comms interface {
	BroadcastRpc(ctx context.Context, msg rpc.Message)
	Rpc(ctx context.Context, id ID, msg rpc.Message)
	Reply() <-chan rpc.Message
}

func ConnectToCluster(cluster *Cluster, virtualLan *LAN) (Comms, error) {
	conns, ok := virtualLan.GetMulticastConns(cluster.SelfID)
	if !ok {
		return nil, fmt.Errorf("can't connect to cluster")
	}

	return NewChannelComms(conns), nil
}

type TCPNetworkComms struct {
}

type ChannelComms struct {
	broadcastOut chan rpc.Message
	rpcChannels  map[ID]chan rpc.Message
	replyChannel chan rpc.Message
	shuttingDown int32
	doneOnce     sync.Once
}

func NewChannelComms(conns map[ID]chan rpc.Message) *ChannelComms {

	comms := &ChannelComms{
		broadcastOut: make(chan rpc.Message),
		rpcChannels:  conns,
		replyChannel: make(chan rpc.Message),
		doneOnce:     sync.Once{},
	}
	return comms
}

func (comms *ChannelComms) Start() {
	fanOut := func(in chan rpc.Message, receivers map[ID]chan rpc.Message) {
		select {
		case msg := <-in:
			for _, r := range receivers {
				select {
				case r <- msg:

				}
			}
		}
	}

	pipeOut := func(c <-chan rpc.Message, out chan rpc.Message) {
		for msg := range c {
			out <- msg
		}
	}

	fanIn := func(upstreams map[ID]chan rpc.Message, out chan rpc.Message) {
		for _, s := range upstreams {
			go pipeOut(s, out)
		}
	}

	go fanOut(comms.broadcastOut, comms.rpcChannels)
	go fanIn(comms.rpcChannels, comms.replyChannel)
}

func (comms *ChannelComms) BroadcastRpc(ctx context.Context, msg rpc.Message) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case comms.broadcastOut <- msg:
				return
			}
		}
	}()
}

func (comms *ChannelComms) Rpc(ctx context.Context, id ID, msg rpc.Message) {
	conn, ok := comms.rpcChannels[id]
	if !ok {
		panic(fmt.Errorf("failed to send rpc: connection with id %d does not exist", id))
	}
	go func() {
		for {
			select {
			case conn <- msg:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (comms *ChannelComms) Reply() <-chan rpc.Message {
	return comms.replyChannel
}

func (comms *ChannelComms) shutdown() {
	atomic.CompareAndSwapInt32(&comms.shuttingDown, 0, 1)
	close(comms.replyChannel)
}

func (comms *ChannelComms) Close() {
	comms.doneOnce.Do(comms.shutdown)
}


