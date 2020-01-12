package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Peer struct {
	ID       uint64 `toml:"id"`
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type PeerList = []Peer

type ServerConfig struct {
	RaftID uint64   `toml:"raft_id"`
	Peers  PeerList `toml:"peers"`

	LocalStoragePath string `toml:"local_storage_path"`

	GrpcServerAddr string `toml:"grpc_server"`
}

func ServerConfigFromTOML(path string) (ServerConfig, error) {
	content, err := ioutil.ReadFile(path)
	conf := ServerConfig{}
	if err != nil {
		return ServerConfig{}, err
	}

	if _, err := toml.Decode(string(content), &conf); err != nil {
		return ServerConfig{}, err
	}

	return conf, nil
}
