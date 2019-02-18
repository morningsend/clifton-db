package blockstore

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
