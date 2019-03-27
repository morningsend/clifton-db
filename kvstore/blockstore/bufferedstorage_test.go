package blockstore

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type testFunc func(storage BlockStorage, t *testing.T)

func runTestInTempFolder(test testFunc, t *testing.T) {
	var testDir, _ = ioutil.TempDir("/tmp", "cliftondb_test")

	defer func() {
		err := os.RemoveAll(testDir)
		if err != nil {
			t.Error(err)
		}
	}()

	storage, err := NewBlockFile(
		path.Join(testDir, "file.dat"),
		WithBlockSize(BaseBlockSize),
	)

	if err != nil {
		t.Error("error creating new block file,", err)
		return
	}

	test(storage, t)
}

func TestBufferedBlockStorage_SequencialReads(t *testing.T) {
	runTestInTempFolder(seqReadTest, t)
}

func TestBufferedBlockStorage_SequencialWrites(t *testing.T) {
	runTestInTempFolder(seqWriteTest, t)
}

func TestBufferedBlockStorage_RandomReads(t *testing.T) {
	runTestInTempFolder(randomReadTest, t)
}

func TestBufferedBlockStorage_RandomWrites(t *testing.T) {
	runTestInTempFolder(randomWriteTest, t)
}
