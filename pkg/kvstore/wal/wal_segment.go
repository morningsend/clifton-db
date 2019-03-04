package wal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
)

type WALSegFlag uint32

const WALSegHeaderMagic = 0x19191919

const (
	WALSegOngoingFlag WALSegFlag = 1 << iota
	WALSegArchivedFlag
	WalSegCompressedFlag
)

type WALSegHeader struct {
	Magic     uint32
	SegId     uint32
	PrevSegId uint32
	Flags     WALSegFlag

	StartRecordIndex uint64
}

type WALSeg struct {
	WALSegHeader

	writeBuffer *bytes.Buffer
	file        *os.File
	logFile     *os.File
	SyncIO      bool
	FilePath    string
	LogSize     int64

	nextRecordIndex uint64
}

type WALRecordReader interface {
	// passes result to WALRecord pointer, avoids allocation
	Read(record *WALRecord) error
	Close() error
}

type WALSegRecordReader struct {
	FilePath string
	file     *os.File

	readBuffer   *bytes.Buffer
	currentIndex uint64
}

func (s *WALSeg) Sync() error {
	if s.logFile == nil {
		return nil
	}

	return s.logFile.Sync()
}

// appending nil record is a no-op
func (s *WALSeg) Append(record *WALRecord) error {
	if record == nil {
		return nil
	}

	record.Index = s.NextRecordIndex()
	s.writeBuffer.Reset()
	record.Marshall(s.writeBuffer)

	var (
		n            int
		err          error
		recordBytes  = s.writeBuffer.Bytes()
		bytesWritten = 0
		sizeToWrite  = s.writeBuffer.Len()
	)

	for bytesWritten < sizeToWrite {
		n, err = s.logFile.Write(recordBytes)
		recordBytes = recordBytes[n:]
		bytesWritten += n
		if err != nil {
			return err
		}
	}

	s.nextRecordIndex = record.Index + 1

	return nil
}

func (s *WALSeg) CreateFileOrLoad(mustNotExist bool) error {
	var err error

	var fileFlags int = os.O_RDWR | os.O_CREATE

	if mustNotExist {
		fileFlags |= os.O_EXCL
	}

	if s.SyncIO {
		fileFlags |= os.O_SYNC
	}

	if s.file == nil {
		s.file, err = os.OpenFile(s.FilePath, fileFlags, 0644)
	}

	if err != nil {
		s.file = nil
		return err
	}

	stat, err := s.file.Stat()

	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		err = s.ReadHeader()
	} else {
		err = s.WriteHeader()
	}

	if err != nil {
		return err
	}

	defer func() {
		err := s.file.Close()
		if err != nil {
			log.Println("error closing file", err)
		}
		s.file = nil
	}()

	return err
}

func (s *WALSeg) ReadHeader() error {
	headerSize := int(unsafe.Sizeof(s.WALSegHeader))
	buffer := make([]byte, headerSize, headerSize)
	bufferTemp := buffer
	bytesRead := 0
	for bytesRead < headerSize {
		n, err := s.file.ReadAt(bufferTemp, int64(bytesRead))
		bytesRead += n
		bufferTemp = bufferTemp[bytesRead:]
		if err != nil {
			return nil
		}
	}

	s.DecodeFromBytes(buffer)

	return nil
}

func NewWALSegment(path string, id uint32, prevId uint32, startIndex uint64, syncIO bool) (*WALSeg, error) {
	seg := &WALSeg{
		WALSegHeader: WALSegHeader{
			Magic: WALSegHeaderMagic,

			SegId:            id,
			PrevSegId:        prevId,
			Flags:            WALSegOngoingFlag,
			StartRecordIndex: startIndex,
		},
		nextRecordIndex: startIndex,
		FilePath:        path,
		writeBuffer:     bytes.NewBuffer(nil),
	}

	err := seg.CreateFileOrLoad(true)

	if err != nil {
		return nil, err
	}

	return seg, nil
}

func OpenWALSegment(path string, syncIO bool) (*WALSeg, error) {
	seg := &WALSeg{
		FilePath:        path,
		SyncIO:          syncIO,
		writeBuffer:     bytes.NewBuffer(nil),
		nextRecordIndex: 0x0,
	}

	err := seg.CreateFileOrLoad(false)

	if err != nil {
		return nil, err
	}

	seg.nextRecordIndex = seg.StartRecordIndex

	if seg.Magic != WALSegHeaderMagic {
		return nil, fmt.Errorf("Segment header magic does not match 0x%x", WALSegHeaderMagic)
	}

	return seg, nil
}

func (s *WALSeg) Archive() error {
	var err error
	if s.file == nil {
		err = s.openForHeaderWriting()
	}

	if err != nil {
		return err
	}

	defer func() {
		err := s.file.Close()
		if err != nil {
			log.Println("error closing file", err)
		}
		s.file = nil
	}()

	s.Flags = WALSegArchivedFlag
	return s.WriteHeader()
}

func (s *WALSeg) PrepareForLogging() error {
	var err error
	if s.logFile == nil {
		s.logFile, err = os.OpenFile(
			s.FilePath,
			os.O_RDWR|os.O_SYNC|os.O_APPEND,
			0644,
		)
	}

	if err != nil {
		return err
	}
	_, err = s.logFile.Seek(0, io.SeekEnd)

	if err != nil {
		err = s.logFile.Close()
		if err != nil {
			log.Println("error closing WAL segment file", err)
		}

		s.logFile = nil
		return err
	}

	return nil
}

func (s *WALSeg) SetReadIndex(logIndex uint64) error {
	if s.StartRecordIndex > logIndex {
		return fmt.Errorf(
			"cannot set read index to be less than start index, start index %d, attempt to set to %d",
			s.StartRecordIndex,
			logIndex,
		)
	}

	return nil
}

func (s *WALSegHeader) EncodeToBytes(buffer []byte) {
	binary.BigEndian.PutUint32(buffer[0:4], s.Magic)
	binary.BigEndian.PutUint32(buffer[4:8], s.SegId)
	binary.BigEndian.PutUint32(buffer[8:12], s.PrevSegId)
	binary.BigEndian.PutUint32(buffer[12:16], uint32(s.Flags))
	binary.BigEndian.PutUint64(buffer[16:24], s.StartRecordIndex)
}

func (s *WALSegHeader) DecodeFromBytes(buffer []byte) {
	s.Magic = binary.BigEndian.Uint32(buffer[0:4])
	s.SegId = binary.BigEndian.Uint32(buffer[4:8])
	s.PrevSegId = binary.BigEndian.Uint32(buffer[8:12])
	s.Flags = WALSegFlag(binary.BigEndian.Uint32(buffer[12:16]))
	s.StartRecordIndex = binary.BigEndian.Uint64(buffer[16:24])
}

func (s *WALSeg) openForHeaderWriting() error {
	var err error
	if s.file == nil {
		s.file, err = os.OpenFile(s.FilePath, os.O_WRONLY, 0644)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *WALSeg) WriteHeader() error {

	buffer := make([]byte, unsafe.Sizeof(s.WALSegHeader))
	s.EncodeToBytes(buffer)
	written := 0
	for written < len(buffer) {
		n, err := s.file.WriteAt(buffer, 0)
		written += n
		if err != nil {
			return err
		}
		buffer = buffer[n:]
	}
	return nil
}

func (s *WALSeg) Close() error {
	var err1, err2 error
	if s.file != nil {
		err1 = s.file.Close()
	}
	if s.logFile != nil {
		err2 = s.logFile.Close()
	}

	if err1 != nil || err2 != nil {
		return fmt.Errorf("error closing files: %v, %v", err1, err2)
	}

	return nil
}

func (s *WALSeg) NextRecordIndex() uint64 {
	return s.nextRecordIndex
}

func (s *WALSeg) NewReader() (WALRecordReader, error) {
	reader := &WALSegRecordReader{
		FilePath:   s.FilePath,
		file:       nil,
		readBuffer: bytes.NewBuffer(nil),
	}

	return reader, nil
}

func (r *WALSegRecordReader) openForReading() error {
	var err error
	r.file, err = os.OpenFile(
		r.FilePath,
		os.O_RDONLY,
		0,
	)

	if err != nil {
		r.file = nil
		return nil
	}

	return nil
}

func (r *WALSegRecordReader) Read(record *WALRecord) error {

	return nil
}

func (r *WALSegRecordReader) Close() error {
	return nil
}
