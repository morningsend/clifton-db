package wal

import (
	"bytes"
	"io"
)

type WALRecordMarshaller interface {
	ToWriter(record *WALRecord, writer io.Writer) error
	ToBuffer(record *WALRecord, buffer *bytes.Buffer) error
}

type WALEventMarshaller interface {
	ToWriter(event *WALEvent, writer io.Writer) error
	ToBuffer(event *WALEvent, buffer *bytes.Buffer) error
}

type WALRecordUnMarshaller interface {
	FromReader(record *WALRecord, reader io.Reader) error
	FromBuffer(record *WALRecord, buffer *bytes.Buffer) error
}

type WALEventUnMarshaller interface {
	FromReader(event *WALEvent, reader io.Reader) error
	FromBuffer(event *WALEvent, buffer *bytes.Buffer) error
}
