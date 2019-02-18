package hash

type Hasher interface {
	WriteString(key string)
	WriteBytes(key []byte)
	GetHash() int64
	GetHashInt32() int32
}

type SimpleHasher int64

func NewSimpleKeyHasher() Hasher {
	return nil
}

func (h *SimpleHasher) WriteString(key string) {

}

func (h *SimpleHasher) WriteBytes(key []byte) {

}

func (h *SimpleHasher) GetHash() int64 {
	return int64(*h)
}

func (h *SimpleHasher) GetHashInt32() int32 {
	return int32(*h)
}
