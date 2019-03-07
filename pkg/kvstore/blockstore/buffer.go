package blockstore

import "bytes"

func clearBuffer(b *bytes.Buffer, value byte, size int) {
	if size > b.Cap() {
		size = b.Cap()
	}
	regionToClear := b.Bytes()[0:size]

	for i, _ := range regionToClear {
		b.Bytes()[i] = value
	}
}
