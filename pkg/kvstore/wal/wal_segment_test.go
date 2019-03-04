package wal

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func getSegmentPath(filename string) string {
	var walDir, _ = ioutil.TempDir("/tmp", "cliftondb_test")
	segmentPath := path.Join(walDir, filename)
	return segmentPath
}

func TestNewWALSegment(t *testing.T) {

	segmentPath := getSegmentPath("segment_00000")
	_ = os.Remove(segmentPath)

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

	defer os.Remove(segmentPath)
}

func TestWALSeg_Archive(t *testing.T) {
	segPath := getSegmentPath("segment_00001")
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
}

func TestWALSeg_Append(t *testing.T) {
	segPath := getSegmentPath("segment_00001")
	const startIndex uint64 = 100
	seg, err := NewWALSegment(segPath, 123, 123, startIndex, false)
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
	
	defer os.Remove(segPath)
}
