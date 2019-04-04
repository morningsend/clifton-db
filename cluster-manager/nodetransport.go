package cluster_manager

import (
	"errors"
	"sync"
)

var NodeIdNonNilErr = errors.New("NodeId must be non nil")
var NodeAlreadyExistsErr = errors.New("node already exists")

const (
	defaultMaxPendingReq = 64
)

type Request struct {
	Id uint64
}

type RequestQueue struct {
	lock sync.RWMutex
}

func NewRequestQueue() *RequestQueue {
	return &RequestQueue{}
}

type ClusterNode struct {
	NodeId         uint32
	NetworkAddress string
}

type ClusterNodeTransportState struct {
	NodeId            uint32
	Client            *ClusterNodeClient
	PendingRequests   *RequestQueue
	MaxPendingRequest int
	NetworkAddress    string
}

type NodeTransport struct {
	nodesLock           sync.RWMutex
	nodes               map[uint32]*ClusterNode
	nodesTransportState map[uint32]*ClusterNodeTransportState
}

func NewNodeTransport() *NodeTransport {
	return &NodeTransport{
		nodes:               make(map[uint32]*ClusterNode),
		nodesTransportState: make(map[uint32]*ClusterNodeTransportState),
	}
}
func (t *NodeTransport) SendToAll() {

}

func (t *NodeTransport) AddNode(node ClusterNode) error {
	t.nodesLock.Lock()
	defer t.nodesLock.Unlock()

	if node.NodeId == 0 {
		return NodeIdNonNilErr
	}

	_, ok := t.nodes[node.NodeId]

	if ok {
		return NodeAlreadyExistsErr
	}

	t.nodes[node.NodeId] = &ClusterNode{
		NodeId:         node.NodeId,
		NetworkAddress: node.NetworkAddress,
	}

	state, err := t.connectToNode(node)

	if err != nil {
		return nil
	}

	t.nodesTransportState[node.NodeId] = state
	return nil
}

func (t *NodeTransport) RemoveNode(nodeId uint32) error {
	t.nodesLock.Lock()
	defer t.nodesLock.Unlock()

	_, ok := t.nodes[nodeId]
	if !ok {
		return nil
	}
}

func (t *NodeTransport) disconnectClient(state *ClusterNodeTransportState) {
	state.Client.Disconnect()
}
func (t *NodeTransport) connectToNode(node ClusterNode) (*ClusterNodeTransportState, error) {
	client, err := NewClusterNodeClient(node.NetworkAddress)
	if err != nil {
		return nil, err
	}

	return &ClusterNodeTransportState{
		NodeId:            node.NodeId,
		Client:            client,
		PendingRequests:   NewRequestQueue(),
		MaxPendingRequest: defaultMaxPendingReq,
		NetworkAddress:    node.NetworkAddress,
	}, nil
}
