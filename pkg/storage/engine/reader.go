package engine

type SliceResult interface {
	Ok() bool
	Err() error
	Value() []byte
}

type KvReader interface {
	Get(k []byte) SliceResult
	MultiGet(k [][]byte) []SliceResult
}
