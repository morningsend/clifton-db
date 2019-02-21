package kvstore

type KVStore interface {
	Get(key string) (data []byte, ok bool, err error)
	Put(key string, data [] byte) (err error)
	Exists(key string) (ok bool, err error)
}
