package wal

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"path/filepath"
)
const WALLockFileName = "wal_lock_file"

type WALCloser interface {
	Close() error
}

type WALSeeker interface {
	ReadIndex() int
	SetIndex(index uint64) error
}

type WALReader interface {
	WALSeeker
	WALCloser
	ReadNext(record *WALRecord) error
}

type WALWriter interface {
	WALCloser
	Append(record *WALRecord) error
	Sync() error
}

type WAL struct {
	DirPath  string
	Segments []WALSeg
	Current  *WALSeg
	AutoSync bool

	CommitIndex uint64
	Index       uint64
}

type walLockFileContent struct {
	CommitIndex uint64 `yaml:"commit-index"`
}

type WALOptions interface {
	Apply(wal *WAL)
}

type autoSyncOptions struct{}

func (*autoSyncOptions) Apply(wal *WAL) {
	wal.AutoSync = true
}

type cleanUpOptions struct{}

func (*cleanUpOptions) Apply(wal *WAL) {

}

func WithAutoSync() WALOptions {
	return &autoSyncOptions{}
}

func WithCleanUp() WALOptions {
	return &cleanUpOptions{}
}

func NewWAL(dirPath string, options ...WALOptions) *WAL {
	wal := &WAL{
		DirPath:     dirPath,
		CommitIndex: 0,

		AutoSync: false,
		Segments: make([]WALSeg, 0, 8),
		Current:  nil,
	}

	for _, opt := range options {
		opt.Apply(wal)
	}

	return wal
}

func (w *WAL) TryRestoreFromLockFile() error {
	lockFilePath := path.Join(w.DirPath, WALLockFileName)
	file, err := os.OpenFile(
		lockFilePath,
		os.O_RDONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	lockFile := walLockFileContent{}
	err = yaml.NewDecoder(file).Decode(&lockFile)

	if err != nil {
		return err
	}

	w.CommitIndex = lockFile.CommitIndex
	return nil
}

func (w *WAL) WriteLockFile() error {
	lockFilePath := path.Join(w.DirPath, WALLockFileName)

	file, err := os.OpenFile(
		lockFilePath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal("error writing lock file at path", lockFilePath)
		}
	}()

	lockFile := walLockFileContent{
		CommitIndex: w.CommitIndex,
	}

	err = yaml.NewEncoder(file).Encode(&lockFile)

	if err != nil {
		return err
	}

	err = file.Sync()

	if err != nil {
		return err
	}

	return nil
}

func RestoreFrom(dirPath string, options ...WALOptions) *WAL {
	wal := &WAL{
		CommitIndex: 0,

		DirPath:  dirPath,
		AutoSync: false,
		Segments: make([]WALSeg, 16),
		Current:  nil,
	}

	for _, opts := range options {
		opts.Apply(wal)
	}

	return wal
}

func (wal *WAL) Sync() error {
	return wal.Current.Sync()
}

func (wal *WAL) NewReader() WALReader {
	return nil
}

func (wal *WAL) NewSegment() error {
	newSeg := WALSeg{
		WALSegHeader: WALSegHeader{
			SegId:     0,
			PrevSegId: 0,

			Magic:            WALSegHeaderMagic,
			Flags:            WALSegOngoingFlag,
			StartRecordIndex: wal.Index,
		},
		file: nil,
	}

	if wal.Current == nil {

		wal.Segments = append(wal.Segments, newSeg)
	}

	return nil
}

func (wal *WAL) LoadSegments() error {
	handleFile := func(path string, in os.FileInfo, err error) error {
		return nil
	}
	err := filepath.Walk(wal.DirPath, handleFile)
	return err
}

func (wal *WAL) Append(record *WALRecord) error {

	if wal.Current == nil {
		seg := WALSeg{}
		wal.Segments = append(wal.Segments, seg)
		wal.Current = &seg
	}

	err := wal.Current.Append(record)

	if err != nil {
		return nil
	}

	if wal.AutoSync {
		err = wal.Current.Sync()
	}

	return err
}
