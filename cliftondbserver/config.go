package cliftondbserver

type RaftNodes struct {
	SelfId   uint32 `yaml:"self-id"`
	Port     uint32 `yaml:"port"`
	PeerList []Peer `yaml:"peers"`
}

type ApiServer struct {
	ListenPort uint32
}

type Peer struct {
	Id       uint32 `yaml:"id"`
	IpOrHost string `yaml:"host"`
	Port     uint32 `yaml:"port"`
}

type Config struct {
	Server ApiServer `yaml:"grpc-grpcServer"`
	DbPath string    `yaml:"db-path"`
	Nodes  RaftNodes `yaml:"raft-nodes"`
}
