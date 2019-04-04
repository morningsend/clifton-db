package cluster_manager

import (
	"github.com/zl14917/MastersProject/api/cluster-services"
	"google.golang.org/grpc"
	"log"
	"net"
)

func init() {

}

func main() {
	grpcServer := grpc.NewServer()
	clusterManagerServer, err := NewClusterManagerService()

	if err != nil {
		log.Fatal("cannot create cluster manager services", err)
	}

	cluster_services.RegisterClusterManagerServer(grpcServer, clusterManagerServer)

	listener, err := net.Listen("tcp", ":9092")

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Printf("server stopped (%v)", err)
	}
}
