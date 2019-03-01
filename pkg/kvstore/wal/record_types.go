package wal

import (
	"bytes"
)

type WALEventType uint32

const (
	PutKey WALEventType = iota
	DeleteKey
)

// Redo logging
type WALEvent struct {
	EventType WALEventType
	KeyLen    uint32
	ValueLen  uint32
	KeyData   []byte
	ValueData []byte
}

type WALRecordHeader struct {
	Index uint32
	CRC   uint32
	Len   uint32
}

type WALRecord struct {
	WALRecordHeader
	WALEvent
}

func (event *WALEvent) Marshall(buffer *bytes.Buffer) {
	//intBytes := make([]byte, 0, 4)
}

func (r *WALRecord) Marshall(buffer *bytes.Buffer) {

}
