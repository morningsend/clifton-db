package wal

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var tmpDir = "/tmp"

type testWithFileIO func(t *testing.T, dirPath string)

func WithTempDir(t *testing.T, io testWithFileIO) {

	name, err := ioutil.TempDir(tmpDir, "cliftondbtests")

	if err != nil {
		t.Fatal("can't create test directory", err)
	}

	io(t, name)

	_ = os.Remove(name)
}

func TestNewWALSegment(t *testing.T) {

	WithTempDir(t, func(t *testing.T, dirPath string) {
		var (
			segmentPath = path.Join(dirPath, "segment_0000")
		)

		seg, err := NewWALSegment(segmentPath, 1000, 999, 1234, true)

		if err != nil {
			t.Error("error creating new wal segment", err)
		}

		if (seg.Flags | WALSegOngoingFlag) < 1 {

		}

		err = seg.Close()

		if err != nil {
			t.Error("should not have error when creating segment file", err)
		}
		if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
			t.Error("should create segment file at", segmentPath, err)
		}

		seg, err = OpenWALSegment(segmentPath, true)

		if err != nil {
			t.Error("failed to reopen segment file", segmentPath, err)
		}

		if seg.SegId != 1000 {
			t.Error("header wrong")
		}

		if seg.PrevSegId != 999 {
			t.Error("header wrong")
		}

		if seg.StartRecordIndex != 1234 {
			t.Error("header wrong")
		}

	})

}

func TestWALSeg_Archive(t *testing.T) {
	WithTempDir(t, func(t *testing.T, dirPath string) {
		segPath := path.Join(dirPath, "segment_00001")
		seg, err := NewWALSegment(segPath, 123, 1222, 1, false)
		//defer seg.Close()
		if err != nil {
			t.Error("error creating segment", err)
		}

		if seg.Flags|WALSegOngoingFlag < 1 {
			t.Error("new segments should be ongoing")
		}

		err = seg.Archive()
		if err != nil {
			t.Error("error archiving segment", err)
		}

		if seg.Flags|WALSegArchivedFlag < 1 {
			t.Error("flags should be WALSegArchivedFlag after archiving.")
		}
	})
}

func TestWALSeg_Append(t *testing.T) {
	WithTempDir(t, func(t *testing.T, dirPath string) {
		segPath := path.Join(dirPath, "segment_00001")

		const startIndex uint64 = 100
		seg, err := NewWALSegment(segPath, 123, 122, startIndex, false)
		if err != nil {
			t.Error("error creating segment", err)
		}

		err = seg.PrepareForLogging()

		if err != nil {
			t.Error("error preparing for logging")
		}

		record := &WALRecord{}
		record.SetPayload(PutKey, []byte("hello"), []byte("world"))
		err = seg.Append(record)

		if err != nil {
			t.Error("error appending record", err)
		}

		if record.Index != startIndex {
			t.Error("appending record should set record index")
		}

	})
}

func TestSegmentReader(t *testing.T) {
	WithTempDir(t, func(t *testing.T, dirPath string) {
		record := &WALRecord{}

		segmentPath := path.Join(dirPath, "segment_0001")

		seg, err := NewWALSegment(segmentPath, 1, 0, 1, false)
		if err != nil {
			t.Error(err)
			return
		}

		reader, err := seg.NewReader()
		if err != nil {
			t.Errorf("error opening reader: %v", err)
		}

		err = reader.Read(record)
		if err != nil {

		}
	})
}
