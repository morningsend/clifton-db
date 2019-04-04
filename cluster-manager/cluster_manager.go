package cluster_manager

import (
	"context"
	"github.com/zl14917/MastersProject/api/cluster-services"
)

type ClusterManagerService struct {
}

//var _ cluster_services.ClusterManagerServer = ClusterManagerService{}

func NewClusterManagerService() (*ClusterManagerService, error) {
	s := &ClusterManagerService{}
	return s, nil
}

func (ClusterManagerService) RegisterNode(context.Context, *cluster_services.RegisterNodeReq) (*cluster_services.RegisterNodeRes, error) {
	return &cluster_services.RegisterNodeRes{}, nil
}

func (ClusterManagerService) DeRegisterNode(context.Context, *cluster_services.DeRegisterNodeReq) (*cluster_services.DeRegisterNodeRes, error) {
	return &cluster_services.DeRegisterNodeRes{}, nil
}
