// Configuration values for KV Store
package kvstore

type SSTableConfig struct {
	IndexBlockSize int `yaml:"index-block-size"`
	DataBlockSize  int `yaml:"data-block-size"`
}

type LogStorageConfig struct {
	SegmentSize string `yaml:"segment-size"`
}

type KVStoreConfig struct {
	Log         LogStorageConfig `yaml:"wal-log"`
	SSTable     SSTableConfig    `yaml:"sstable"`
	DataDirPath string           `yaml:"data-dir"`
}
