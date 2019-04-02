package kvserver

import "github.com/zl14917/MastersProject/kvstore"

type Cluster struct {
	SelfId   uint32 `yaml:"self-id"`
	Port     uint32 `yaml:"port"`
	PeerList []Peer `yaml:"peers"`
}

type Peer struct {
	Id       uint32 `yaml:"id"`
	IpOrHost string `yaml:"host"`
	Port     uint32 `yaml:"port"`
}

type Config struct {
	Data   kvstore.KVStoreConfig `yaml:"kv-store"`
	DbPath string                `yaml:"db-path"`
	Nodes  Cluster               `yaml:"nodes"`
}
