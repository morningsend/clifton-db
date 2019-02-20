package raft

import (
	"context"
	"github.com/zl14917/MastersProject/pkg/raft/rpc"
	"log"
	"os"
)

type Role int
type Status int
type ID int

const (
	Follower Role = iota
	NoneVoteFollower
	Candidate
	Leader
)

const (
	NoTimeOut = iota
	ElectionTimedOut
	HeartBeatTimeOut
)

type CommonState struct {
	CurrentTerm int
	VotedFor    ID

	CommitIndex      int
	LastAppliedIndex int
}

func (s *CommonState) NextTerm() {
	s.CurrentTerm++
}

type CandidateState struct {
	Majority         int
	ElectionDeadline int
	VotesReceived    map[ID]bool
}

func (c *CandidateState) ResetVotes() {
	for k, _ := range c.VotesReceived {
		c.VotesReceived[k] = false
	}
}

func (s *CandidateState) SetVote(id ID) {
	s.VotesReceived[id] = true
}

func (s *CandidateState) HasMajority() bool {
	return s.CountVote() >= s.Majority
}

func (c *CandidateState) CountVote() int {
	count := 0
	for _, v := range c.VotesReceived {
		if v {
			count++
		}
	}
	return count
}

type LeaderState struct {
	NextIndices  map[ID]int
	MatchIndices map[ID]int
}

func (s *LeaderState) Reset() {
	for k, _ := range s.NextIndices {
		s.NextIndices[k] = 0
	}

	for k, _ := range s.MatchIndices {
		s.MatchIndices[k] = 0
	}
}

type FollowerState struct {
	HeartBeatDeadline int
}

type RaftFSM interface {
	raftFSM()
	Id() ID
	Role() Role
	GetCurrentTerm() int
	CommitIndex() int
	LastAppliedIndex() int
	VotedFor() ID
	Log() RaftLog

	Tick() int
}

type RealRaftFSM struct {
	CommonState
	LeaderState
	FollowerState
	CandidateState

	role             Role
	id               ID
	currentTick      int
	electionTimeout  int
	heartBeatTimeout int
	comms            Comms
	logger           *log.Logger
	peers            []ID
	log              RaftLog
}

func NewRaftFSM(id ID, peers []ID, electionTimeout int, heartBeatTimeout int, comms Comms) RaftFSM {
	fsm := RealRaftFSM{
		id:               id,
		role:             Follower,
		currentTick:      0,
		electionTimeout:  electionTimeout,
		heartBeatTimeout: heartBeatTimeout,
		peers:            peers,
		log:              NewInMemoryLog(),
		logger:           log.New(os.Stdout, "[raft]", log.LstdFlags),
		comms:            comms,
		CandidateState: CandidateState{
			VotesReceived: make(map[ID]bool),
		},
	}
	fsm.BecomeFollower()

	return &fsm
}

func (r *RealRaftFSM) raftFSM() {}

func (r *RealRaftFSM) Role() Role {
	return r.role
}

func (r *RealRaftFSM) Id() ID {
	return r.id
}

func (r *RealRaftFSM) GetCurrentTerm() int {
	return r.CommonState.CurrentTerm
}

func (r *RealRaftFSM) VotedFor() ID {
	return r.id
}

func (r *RealRaftFSM) CommitIndex() int {
	return r.CommonState.CommitIndex
}

func (r *RealRaftFSM) LastAppliedIndex() int {
	return r.CommonState.LastAppliedIndex
}

func (r *RealRaftFSM) Log() RaftLog {
	return r.log
}

func (r *RealRaftFSM) tick() (int) {
	r.currentTick++
	r.checkTimeout()
	r.step()
	return r.currentTick
}

func (r *RealRaftFSM) Tick() int {
	return r.tick()
}
func (r *RealRaftFSM) checkTimeout() {
	switch r.role {
	case Follower:
		if r.currentTick >= r.FollowerState.HeartBeatDeadline {
			r.BecomeCandidate()
			return
		}
	case Candidate:
		if r.currentTick >= r.CandidateState.ElectionDeadline {
			r.BecomeFollower()
			return
		}
	case Leader:
		return
	default:
		return
	}
}

func (r *RealRaftFSM) step() {
	switch r.role {
	case NoneVoteFollower:
		return
	case Follower:
		r.followerStep()
	case Candidate:
		r.candidateStep()
	case Leader:
		r.leaderStep()
	}
}
func (r *RealRaftFSM) BecomeFollower() {
	r.role = Follower
	r.FollowerState.HeartBeatDeadline = r.currentTick + r.heartBeatTimeout
}

func (r *RealRaftFSM) BecomeCandidate() {
	r.CandidateState.ElectionDeadline = r.currentTick + r.electionTimeout
	r.role = Candidate
	r.CurrentTerm++
	r.ResetVotes()
	r.SetVote(r.id)

	index, term := r.log.GetLastLogTermIndex()
	requestVoteMsg := &rpc.RequestVote{
		CandidateId:  int(r.id), // TODO: refactor ID type to its own package,
		Term:         r.GetCurrentTerm(),
		LastLogIndex: index,
		LastLogTerm:  term,
	}

	r.broadcastRpcImmediate(requestVoteMsg)
}

func (r *RealRaftFSM) BecomeLeader() {
	r.role = Leader
	r.HeartBeatDeadline = r.currentTick + r.heartBeatTimeout - 2

	// send initial heart beat to all other nodes.
	msg := rpc.AppendEntriesReq{
		Term:              r.CommonState.CurrentTerm,
		LeaderId:          int(r.id),
		LeaderCommitIndex: r.CommonState.CommitIndex,
		PrevLogIndex:      r.CommonState.LastAppliedIndex,
		Entries:           nil,
	}
	r.broadcastRpcImmediate(&msg)
}

func (r *RealRaftFSM) ReceiveRpc(msg rpc.Message) {
	switch x := msg.(type) {
	case *rpc.AppendEntriesReq:
		r.ReceiveAppendEntriesRpc(x)
	case *rpc.AppendEntriesReply:
		r.ReceiveAppendEntriesReply(x)
	case *rpc.RequestVote:
		r.ReceiveVoteRequestRpc(x)
	case *rpc.VotedFor:
		r.ReceiveVotedFor(x)
		return
	default:
		r.logger.Println("error, unrecognized message type:", x)
	}
}

func (r *RealRaftFSM) ReceiveAppendEntriesReply(msg *rpc.AppendEntriesReply) {

}

func (r *RealRaftFSM) ReceiveVoteRequestRpc(msg *rpc.RequestVote) {

}

func (r *RealRaftFSM) ReceiveHeartBeat(msg *rpc.HeartBeat) {

}

func (r *RealRaftFSM) ReceiveAppendEntriesRpc(msg *rpc.AppendEntriesReq) {
	senderId := msg.LeaderId
	if msg.Term > r.CommonState.CurrentTerm {
		r.CommonState.CurrentTerm = msg.Term
		r.BecomeFollower()
	}
	reply := rpc.AppendEntriesReply{
		Term:    r.CommonState.CurrentTerm,
		Success: true,
	}
	r.sendRpcImmediate(ID(senderId), &reply)
	for _, e := range msg.Entries {
		err := r.log.Append(Entry{
			Index: msg.PrevLogIndex,
			Term:  msg.Term,
			Data:  e,
		})

		if err != nil {
			r.logger.Fatalln("cannot append log:", err)
		}
	}
}
func (r *RealRaftFSM) ReceiveVotedFor(msg *rpc.VotedFor) {
	if r.role != Candidate {
		r.unhandledRpc(msg)
		return
	}
}

func (r *RealRaftFSM) unhandledRpc(msg rpc.Message) {

}

func (r *RealRaftFSM) candidateStep() {
	votes := len(r.CandidateState.VotesReceived)
	if votes >= r.CandidateState.Majority {
		r.BecomeLeader()
	}
}

func (r *RealRaftFSM) leaderStep() {
	if r.HeartBeatDeadline <= r.currentTick {
		r.HeartBeatDeadline = r.currentTick + r.heartBeatTimeout - 2
	}
	msg := &rpc.AppendEntriesReq{
		Term:              r.CommonState.CurrentTerm,
		LeaderId:          int(r.id),
		LeaderCommitIndex: r.CommonState.CommitIndex,
		PrevLogIndex:      r.CommonState.LastAppliedIndex,
		Entries:           nil,
	}
	r.broadcastRpcImmediate(msg)
}

func (r *RealRaftFSM) followerStep() {
	// 1. check heartbeat timeout

}

func (r *RealRaftFSM) enqueueBroadcastRpc(msg rpc.Message) {
	go r.comms.BroadcastRpc(context.Background(), msg)
}

func (r *RealRaftFSM) broadcastRpcImmediate(msg rpc.Message) {
	go r.comms.BroadcastRpc(context.Background(), msg)
}

func (r *RealRaftFSM) sendRpcImmediate(receiver ID, msg rpc.Message) {
	go r.comms.Rpc(context.Background(), receiver, msg)
}

func (r *RealRaftFSM) enqueueRpc(receiver ID, msg rpc.Message) {
	go r.comms.Rpc(context.Background(), receiver, msg)
}
