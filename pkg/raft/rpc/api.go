package rpc

type SendType int

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
	Entries           []interface{}
}

type AppendEntriesReply struct {
	Message
	Term int
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
	Term        int
	VoteGranted bool
}
