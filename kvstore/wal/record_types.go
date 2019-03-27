package wal

import (
	"encoding/binary"
	"io"
)

type WALEventType uint32

const (
	PutKey WALEventType = iota
	DeleteKey
)

// Redo logging
type WALEvent struct {
	EventData []byte
}

type WALRecordHeader struct {
	Index     uint64
	CRC       uint32
	DataLen   uint32
	EventType WALEventType
}

type WALRecord struct {
	WALRecordHeader
	WALEvent
}

func (header *WALRecordHeader) Marshall(writer io.Writer) error {
	var (
		err          error
		buffer8Bytes [8]byte
		uint64buf    = buffer8Bytes[0:8]
		uint32buf    = buffer8Bytes[0:4]
	)

	binary.BigEndian.PutUint64(uint64buf, header.Index)
	_, err = writer.Write(uint64buf)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buf, header.CRC)
	_, err = writer.Write(uint32buf)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buf, header.DataLen)
	_, err = writer.Write(uint32buf)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buf, uint32(header.EventType))
	_, err = writer.Write(uint32buf)
	if err != nil {
		return err
	}

	return nil
}

func (h *WALRecordHeader) UnMarshall(r io.Reader) error {
	var (
		err         error
		buffer8Byte [8]byte
		uint32buf   = buffer8Byte[0:4]
		uint64buf   = buffer8Byte[0:8]
	)

	_, err = r.Read(uint64buf)
	if err != nil {
		return err
	}

	h.Index = binary.BigEndian.Uint64(uint64buf)

	_, err = r.Read(uint32buf)
	if err != nil {
		return err
	}

	h.CRC = binary.BigEndian.Uint32(uint32buf)

	_, err = r.Read(uint32buf)
	if err != nil {
		return err
	}
	h.DataLen = binary.BigEndian.Uint32(uint32buf)

	_, err = r.Read(uint32buf)
	if err != nil {
		return err
	}

	h.EventType = WALEventType(binary.BigEndian.Uint32(uint32buf))
	return nil
}

func (r *WALRecord) Marshall(writer io.Writer) error {
	var err error

	err = r.WALRecordHeader.Marshall(writer)
	if err != nil {
		return err
	}

	_, err = writer.Write(r.EventData)
	if err != nil {
		return err
	}

	return nil
}

func ComputeWALRecordCRC(data []byte) uint32 {
	return 0
}

func (w *WALRecord) VerifyCRC() bool {
	crc := ComputeWALRecordCRC(w.EventData)
	return w.CRC == crc
}

func (w *WALRecord) ComputeCRC() {
	if len(w.EventData) == 0 {
		w.CRC = 0
		return
	}
	w.CRC = ComputeWALRecordCRC(w.EventData)
}
