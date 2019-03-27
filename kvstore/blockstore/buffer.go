package blockstore

import "bytes"

func clearBuffer(b *bytes.Buffer, value byte, size int) {
	if size > b.Cap() {
		size = b.Cap()
	}
	regionToClear := b.Bytes()[0:size]
	clearBytes(regionToClear, value, size)
}

func clearBytes(bArray []byte, value byte, size int) {
	for i, _ := range bArray {
		bArray[i] = value
	}
}
