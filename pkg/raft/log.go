package raft

import (
	"fmt"
	"sync"
)

type Entry struct {
	Index int
	Term  int
	Data  interface{}
}

type RaftLog interface {
	GetLastApplied() int
	GetLastCommitted() int
	GetLastLogTermIndex() (term, index int)
	RaftLogAppender
	RaftLogTruncator
	RaftLogReader
	RaftLogCommitter
}

type RaftLogAppender interface {
	Append(entry Entry) error
}

type RaftLogTruncator interface {
	TruncateTo(index int) error
}

type RaftLogReader interface {
	Read(index int) (entry *Entry, ok bool)
}

type RaftLogCommitter interface {
	Commit(index int) error
}

func NewInMemoryLog() RaftLog {
	entries := make([]Entry, 0, 1000)
	//insert sentinel value
	entries = append(entries, Entry{
		0, 0, nil,
	})
	return &InMemoryLog{
		Entries: entries,
		lock:    &sync.RWMutex{},
	}
}

type InMemoryLog struct {
	Entries       []Entry
	LastApplied   uint64
	LastCommitted uint64
	lock          *sync.RWMutex
}

func (l *InMemoryLog) GetLastApplied() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return int(l.LastApplied)
}

func (l *InMemoryLog) GetLastCommitted() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return int(l.LastCommitted)
}

func (l *InMemoryLog) Append(entry Entry) (error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	entry.Index = len(l.Entries)
	l.Entries = append(l.Entries, entry)
	return nil
}

func (l *InMemoryLog) TruncateTo(index int) (error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if len(l.Entries) <= index+1 {
		return nil
	}
	l.Entries = l.Entries[0 : index+1]
	return nil
}

func (l *InMemoryLog) GetLastLogTermIndex() (term, index int) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	lastEntryIndex := len(l.Entries) - 1
	lastEntry := &l.Entries[lastEntryIndex]
	return lastEntry.Term, lastEntry.Index
}

func (l *InMemoryLog) Read(index int) (entry *Entry, ok bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if len(l.Entries) <= index {
		return nil, false
	}

	return &l.Entries[index], true
}

func (l *InMemoryLog) Commit(index int) (error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if (l.LastCommitted + 1) != uint64(index) {
		return fmt.Errorf("commit index is not consecutive from last commited.")
	}
	l.LastCommitted++
	return nil
}
