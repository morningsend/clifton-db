package sstable

import (
	"encoding/binary"
	"errors"
	"github.com/zl14917/MastersProject/pkg/kvstore/types"
	"io"
	"unsafe"
)

const (
	DataFileMagic uint32 = 0x33333333

	SSTableDataFileHeaderSize   = unsafe.Sizeof(SSTableDataFileHeader{})
	SSTableDataRecordHeaderSize = unsafe.Sizeof(SSTableDataRecord{}.Value)
)

type SSTableDataFileHeader struct {
	Magic      uint32
	Flags      DataFileFlags
	BlockSize  uint32
	BlockCount uint32
}

var UnitializedSSTableDataFileHeader = SSTableDataFileHeader{
	Magic:      DataFileMagic,
	Flags:      DataFileFlags(HeaderUninitialized),
	BlockSize:  HeaderUninitialized,
	BlockCount: HeaderUninitialized,
}

func (h *SSTableDataFileHeader) Marshall(w io.Writer) error {
	var (
		buffer       [4]byte
		uint32buffer = buffer[0:4]
		err          error
	)

	binary.BigEndian.PutUint32(uint32buffer, h.Magic)
	_, err = w.Write(uint32buffer)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buffer, uint32(h.Flags))
	_, err = w.Write(uint32buffer)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buffer, h.BlockSize)
	_, err = w.Write(uint32buffer)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(uint32buffer, h.BlockCount)
	_, err = w.Write(uint32buffer)
	if err != nil {
		return err
	}

	return nil
}

func (h *SSTableDataFileHeader) UnMarshall(r io.Reader) error {
	var (
		buffer       [4]byte
		uint32buffer = buffer[0:4]
		err          error
	)

	_, err = r.Read(uint32buffer)
	if err != nil {
		return err
	}

	h.Magic = binary.BigEndian.Uint32(uint32buffer)

	_, err = r.Read(uint32buffer)
	if err != nil {
		return err
	}

	h.Flags = DataFileFlags(binary.BigEndian.Uint32(uint32buffer))

	_, err = r.Read(uint32buffer)
	if err != nil {
		return err
	}

	h.BlockSize = binary.BigEndian.Uint32(uint32buffer)

	_, err = r.Read(uint32buffer)
	if err != nil {
		return err
	}

	h.BlockCount = binary.BigEndian.Uint32(uint32buffer)

	return nil
}

type SSTableDataRecord struct {
	ValueLen uint32
	Value    types.ValueType
}

func MaxValueSizeFitInBlock(blockSize int) int {
	return blockSize - int(SSTableDataRecordHeaderSize)
}

func (d *SSTableDataRecord) Marshall(w io.Writer) error {
	var (
		err       error
		buffer    [4]byte
		uint32buf = buffer[0:4]
	)

	if d.ValueLen != uint32(len(d.Value)) {
		return errors.New("length mismatch error")
	}

	padTo := NextMultipleOf4Uint(uint(d.ValueLen))
	gap := padTo - uint(d.ValueLen)
	binary.BigEndian.PutUint32(uint32buf, d.ValueLen)
	_, err = w.Write(uint32buf)
	if err != nil {
		return err
	}

	if gap > 0 {
		binary.BigEndian.PutUint32(uint32buf, 0)
		_, err = w.Write(uint32buf[0:gap])
	}

	return nil
}

func (d *SSTableDataRecord) UnMarshall(r io.Reader) error {
	var (
		err       error
		buffer    [4]byte
		uint32buf = buffer[0:4]
	)

	_, err = r.Read(uint32buf)
	if err != nil {
		return err
	}

	d.ValueLen = binary.BigEndian.Uint32(uint32buf)
	multipleOf4 := NextMultipleOf4Uint(uint(d.ValueLen))

	d.Value = make([]byte, multipleOf4, multipleOf4)
	_, err = r.Read(d.Value)
	d.Value = d.Value[0:d.ValueLen]

	return err
}
