package engine

type WriteBatch struct {
}

type KvWriter interface {
	Put(k []byte, value []byte) error
	Delete(k []byte) error
}

type KvBatchWriter interface {
	Do(b WriteBatch) error
}
