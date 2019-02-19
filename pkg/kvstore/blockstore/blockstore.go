package blockstore

type BlockStorage interface {
	ReadBytes(pointer BlockPointer) ([]byte, error)
	WriteBytes(buffer []byte, nbytes int) (BlockPointer, error)
}

type BlockPointer struct {
	BlockNumber int32
	Offset      int32
	NumBytes    int32
}

type Block struct {
	BlockHeader
	BlockBody
	Next *Block
}

type BlockHeader struct {
	Length     int32
	BlockSize  int32
	BlockCRC   int64
	Key        []byte
	HeaderData []byte
}

type BlockBody struct {
	Data []byte
}

type InMemoryBlockStorage struct {
	BlockSize int
}

func (store *InMemoryBlockStorage) ReadBytes(pointer BlockPointer) ([]byte, error) {
	return []byte("hello"), nil
}

func (store *InMemoryBlockStorage) WriteBytes(buffer []byte, nbytes int) (BlockPointer, error) {
	return BlockPointer{}, nil
}

func NewInMemoryStore(blockSize int) BlockStorage {
	return &InMemoryBlockStorage{
		BlockSize: blockSize,
	}
}
