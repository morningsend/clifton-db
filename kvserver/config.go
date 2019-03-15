package kvserver

type KVStorageConfig struct {
	IndexBlockSize int `yaml:"index-block-size"`
	DataBlockSize  int `yaml:"data-block-size"`
}

type LogStorageConfig struct {
	SegmentSize string `yaml:"segment-size"`
}

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
	Data   KVStorageConfig  `yaml:"data"`
	Log    LogStorageConfig `yaml:"log"`
	DbPath string           `yaml:"db-path"`
	Nodes  Cluster          `yaml:"nodes"`
}
