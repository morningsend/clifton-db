package cluster_manager

import (
	"github.com/zl14917/MastersProject/api/cluster-services"
	"google.golang.org/grpc"
)

type ClusterNodeClient struct {
	NetworkAddress string
	*grpc.ClientConn
	cluster_services.ClusterNodeClient
}

func NewClusterNodeClient(address string) (*ClusterNodeClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	client := cluster_services.NewClusterNodeClient(conn)

	return &ClusterNodeClient{
		NetworkAddress:    address,
		ClientConn:        conn,
		ClusterNodeClient: client,
	}, nil
}

func (c *ClusterNodeClient) Disconnect() error {
	return c.ClientConn.Close()
}
