package kvstore

type KVStoreReadOptions interface {
}

type KVStore interface {
	Get(key string) (data []byte, ok bool, err error)
	Put(key string, data [] byte) (err error)
	Delete(key string) (ok bool, err error)
	Exists(key string) (ok bool, err error)
}

type KVStoreOpenOptions interface {
	Apply(options *KVStoreOptions)
}

type KVStoreOptions struct {
	WALSegmentSizeBytes uint32
	MaxKeySizeBytes     uint32
	MaxValueSizeBytes   uint32
}

var defaultKVStoreOptions = KVStoreOptions{
	WALSegmentSizeBytes: 1024 * 1024 * 16,
	MaxKeySizeBytes:     1024 * 4,
	MaxValueSizeBytes:   16 * 1024,
}

func OpenKVStore(dirPath string, options ...KVStoreOpenOptions) KVStore {
	var kvStoreConfig = defaultKVStoreOptions

	for _, opts := range options {
		opts.Apply(&kvStoreConfig)
	}

}
