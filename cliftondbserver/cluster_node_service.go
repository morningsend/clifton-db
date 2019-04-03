package cliftondbserver

import (
	"context"
	"github.com/zl14917/MastersProject/api/cluster-node"
	"github.com/zl14917/MastersProject/raft"
)

var _ cluster_node.ClusterNodeServer = GrpcClusterNodeService{}

type GrpcClusterNodeService struct {
	raftNode *raft.RaftNode
}

func (GrpcClusterNodeService) AddNode(context.Context, *cluster_node.AddNodeReq) (*cluster_node.AddNodeRes, error) {
	panic("implement me")
}

func (GrpcClusterNodeService) GetNodeList(context.Context, *interface{}) (*cluster_node.PeerListRes, error) {
	panic("")
}

func (GrpcClusterNodeService) RemoveNode(context.Context, *cluster_node.RemoveNodeReq) (*cluster_node.RemoveNodeRes, error) {
	panic("implement me")
}

func NewClusterNodeService(raft *raft.RaftNode) {

}
