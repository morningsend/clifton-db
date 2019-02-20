package rpc

type SendType int
type Role int
type Status int
type ID int

const (
	Follower Role = iota
	NoneVoteFollower
	Candidate
	Leader
)

type Message interface {
	message()
}

type HeartBeat interface {
	Message
	heartBeat()
}

type AppendEntriesReq struct {
	Message
	Term              int
	LeaderId          int
	LeaderCommitIndex int
	PrevLogIndex      int
	PrevLogTerm       int
	Entries           []interface{}
}

type AppendEntriesReply struct {
	Message
	Term    int
	Success bool
}

type AppendEntriesReplyData struct {
}

type RequestVote struct {
	Message
	CandidateId  int
	Term         int
	LastLogIndex int
	LastLogTerm  int
}

type VotedFor struct {
	Message
	VoterId     int
	Term        int
	VoteGranted bool
}
