package engine

type SliceResult interface {
	Ok() bool
	Err() error
	Value() []byte
}

type KvReader interface {
	Get(k []byte) SliceResult
	MultiGet(k [][]byte) []SliceResult
	Scan(k []byte, prefixIsKey bool) Iterator
}

type Iterator interface {
	Item() (k []byte, value []byte)
	Valid() bool
	Next()
}
