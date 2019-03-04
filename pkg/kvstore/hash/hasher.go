package hash

type SimpleHasher int32

type Hasher interface {
	UpdateString(key string)
	UpdateBytes(key []byte)
	GetValue() int32
}

func NewSimpleKeyHasher() Hasher {
	return nil
}

func (h *SimpleHasher) WriteString(key string) {

}

func (h *SimpleHasher) WriteBytes(key []byte) {

}

func (h *SimpleHasher) GetHash() SimpleHasher {
	return int32(*h)
}
