package wal

import (
	"bytes"
	"io"
	"unsafe"
)

type WALEventType uint32

const (
	PutKey WALEventType = iota
	DeleteKey
)

type WALEventHeader struct {
	EventType WALEventType
	KeyLen    uint32
	ValueLen  uint32
}

// Redo logging
type WALEvent struct {
	WALEventHeader
	KeyData   []byte
	ValueData []byte
}

type WALRecordHeader struct {
	Index uint64
	CRC   uint32
	Len   uint32
}

type WALRecord struct {
	WALRecordHeader
	WALEvent
}

func (event *WALEvent) Marshall(writer io.Writer) {
	//intBytes := make([]byte, 0, 4)

}

func (header *WALRecordHeader) Marshall(buffer *bytes.Buffer) {

}

func (r *WALRecord) Marshall(buffer *bytes.Buffer) {
	r.WALRecordHeader.Marshall(buffer)
}

func (r *WALRecord) SetPayload(eventType WALEventType, key []byte, value []byte) {
	r.EventType = eventType
	r.KeyLen = uint32(len(key))
	r.KeyData = key
	r.ValueData = value
	r.ValueLen = uint32(len(value))
	r.Len = r.KeyLen + r.ValueLen + 2*uint32(unsafe.Sizeof(r.ValueLen))
}

func ComputeWALRecordCRC(data []byte, len int) int32 {

	return 0
}
