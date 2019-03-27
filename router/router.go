package router

import (
	"github.com/zl14917/MastersProject/dht"
	"google.golang.org/grpc"
)

type NodeId int
type NetworkAddress string

type RoutingTable struct {
	NodeNetworkAddresses map[NodeId]NetworkAddress
	HashRing             dht.HashRing
	SelfId               NodeId
}

type Router interface {
}

type ClientRequestRouter struct {
	RoutingTable
	Clients map[NodeId]*grpc.ClientConn
}

func CreateRouter(selfId NodeId) *ClientRequestRouter {
	return &ClientRequestRouter{
		RoutingTable: RoutingTable{
			HashRing:             dht.NewHashRing(),
			SelfId:               selfId,
			NodeNetworkAddresses: make(map[NodeId]NetworkAddress),
		},
	}
}

func (r *ClientRequestRouter) Connect() error {
	return nil
}

func (r *ClientRequestRouter) IsSelf(id NodeId) bool {
	return true
}
