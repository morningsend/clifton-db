package blockstore

import "os"

// A class to provide block buffering for io.
// Data is aligned in blocks
type BufferedBlockStorage struct {
	FilePath string
	fileIO   *os.File

	blockCount uint
	autoSync   bool
}

func (b *BufferedBlockStorage) blockStorage() {

}

func (b *BufferedBlockStorage) Write(p []byte) (n int, err error) {
	return
}

func (b *BufferedBlockStorage) Read(p []byte) (n int, err error) {
	return
}

func NewBlockFile(path string, blockSize uint) (BlockStorage, error) {
	return nil, nil
}

func OpenBlockFile(path string) (BlockStorage, error) {
	return nil, nil
}

