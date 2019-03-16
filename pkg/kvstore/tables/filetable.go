package tables

import (
	"fmt"
	"github.com/zl14917/MastersProject/pkg/kvstore/sstable"
	"strconv"
	"sync"
	"time"
)

type FileTableScanner interface {
}

type MemTableFlushCallback func(bool, error)

type FileTable interface {
	BeginFlushing(table MemTable, withCallback MemTableFlushCallback)
	NewScanner() FileTableScanner
}

type SStableRef struct {
	sync.Mutex
	sstable.SSTable

	Timestamp int64
	Level     int
}

// Leveled SStables are stored in concurrent lists.
//
type LevelFileTable struct {
	TableRootDir string

	Level0 []*SStableRef
	Level1 []*SStableRef
	Level2 []*SStableRef
}

type SSTableFlushEvent struct {
	Level     int
	Timestamp int64
}

var _ FileTable = NewSStableFileTable("", "")

func NewSStableFileTable(tableRootDir string, logDir string) *LevelFileTable {
	table := &LevelFileTable{
		TableRootDir: tableRootDir,
	}
	return table
}

func NewSSTableRef(level int) *SStableRef {
	if level < 0 {
		level = 0
	}

	tablet := &SStableRef{}
	options := sstable.SSTableOpenOptions{
		Prefix:       "level_" + strconv.Itoa(level),
		MaxKeySize:   0,
		MaxValueSize: 0,
		BlockSize:    0,
		KeyBlockSize: int(sstable.BaseBlockSize),
		LoadExisting: false,
		Timestamp:    time.Now().Unix(),
	}
	tablet.SSTable = *(sstable.NewSSTable("", &options))
	return tablet
}

// Flushing Memtable to File Table creates a level 0 SSTable tablet
//
func (t *LevelFileTable) BeginFlushing(table MemTable, withCallback MemTableFlushCallback) {
	var err error

	filename := fmt.Sprintf("level_0_%d")
	fmt.Printf(filename)

	iterator := table.Iterator()

	newSStable := NewSSTableRef(0)
	writer, err := newSStable.NewWriter()

	if err != nil {
		withCallback(false, err)
		return
	}

	for iterator.Next() {

		key, value := iterator.Current()
		err = writer.Write(key, []byte(value), value == nil)

		if err != nil {
			withCallback(false, err)
		}
	}

	event := SSTableFlushEvent{
		Timestamp: newSStable.Timestamp,
		Level:     newSStable.Level,
	}

	err = t.CommitChangeToLockFile(event)

	if err != nil {
		withCallback(false, err)
	}

	withCallback(true, nil)
}

func (t *LevelFileTable) CommitChangeToLockFile(v interface{}) error {
	return nil
}

func (t *LevelFileTable) NewScanner() FileTableScanner {
	return nil
}
