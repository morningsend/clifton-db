package kvserver

type KVStorageConfig struct {
	PathDir string `yaml:"path-dir"`
}

type LogStorageConfig struct {
	PathDir     string `yaml:"path-dir"`
	SegmentSize string `yaml:"segment-size"`
}

type Nodes struct {
	SelfId   uint32           `yaml:"self-id"`
	PeersMap []NetworkAddress `yaml:"peers"`
}

type NetworkAddress struct {
	Id       uint32 `yaml:"id"`
	IpOrHost string `yaml:"host"`
	Port     uint32 `yaml:"port"`
}

type Config struct {
	Data  KVStorageConfig  `yaml:"data"`
	Log   LogStorageConfig `yaml:"log"`
	Nodes Nodes            `yaml:"nodes"`
}
