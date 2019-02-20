package raft

type Entry struct {
	Index int
	Term  int
	Data  interface{}
}

type RaftLog interface {
	GetLastApplied() int
	GetLastCommitted() int

	RaftLogAppender
	RaftLogTruncator
}

type RaftLogAppender interface {
	Append(entry Entry) error
}

type RaftLogTruncator interface {
	TruncateTo(index int) error
}

func NewInMemoryLog() RaftLog {
	return &InMemoryLog{
		Entries: make([]Entry, 0, 1000),
	}
}

type InMemoryLog struct {
	Entries       []Entry
	LastApplied   uint64
	LastCommitted uint64
}

func (l *InMemoryLog) GetLastApplied() int {
	return int(l.LastApplied)
}

func (l *InMemoryLog) GetLastCommitted() int {
	return int(l.LastCommitted)
}

func (l *InMemoryLog) Append(entries Entry) (error) {
	return nil
}

func (l *InMemoryLog) TruncateTo(index int) (error) {
	return nil
}
