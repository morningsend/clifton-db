package raft

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/pkg/raft/rpc"
	"log"
	"os"
	"reflect"
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

type RaftFSM interface {
	raftFSM()
	Id() ID
	Role() Role
	GetCurrentTerm() int
	CommitIndex() int
	LastAppliedIndex() int
	VotedFor() ID
	Log() RaftLog
	GetComms() Comms
	Tick() int
	ReceiveMsg(msg rpc.Message)
	ReplicateDataToLog(data interface{})
}

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
	ElectionTimeout  int
	VotesReceived    map[ID]bool
}

func (c *CandidateState) ResetTime(currentTick int) {
	c.ElectionDeadline = currentTick + c.ElectionTimeout
}

func (c *CandidateState) HasElectionTimeout(currentTick int) bool {
	return currentTick < c.ElectionDeadline
}

func (c *CandidateState) ResetVotes() {
	for k, _ := range c.VotesReceived {
		c.VotesReceived[k] = false
	}
}

func (s *CandidateState) GotVoteFrom(id ID) {
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
	NextIndices       map[ID]int
	MatchIndices      map[ID]int
	NextHeartBeatTime int
	HeartBeatTimeout  int
}

func (s *LeaderState) Reset() {
	for k, _ := range s.NextIndices {
		s.NextIndices[k] = 0
	}

	for k, _ := range s.MatchIndices {
		s.MatchIndices[k] = 0
	}
}

func (s *LeaderState) ResetNextHeartBeatTime(currentTick int) {
	s.NextHeartBeatTime = currentTick + s.HeartBeatTimeout - 3
}

type FollowerState struct {
	HeartBeatDeadline int
	HeartBeatTimeout  int
}

func (s *FollowerState) ResetHeartBeatTimeout(currentTick int) {
	s.HeartBeatDeadline = currentTick + s.HeartBeatTimeout
}

func (s *FollowerState) HasHeartBeatTimeout(currentTick int) bool {
	return currentTick < s.HeartBeatDeadline
}

type RealRaftFSM struct {
	CommonState
	LeaderState
	FollowerState
	CandidateState

	role        Role
	id          ID
	currentTick int
	comms       Comms
	logger      *log.Logger
	peers       []ID
	log         RaftLog
}

func NewRaftFSM(id ID, peers []ID, electionTimeout int, heartBeatTimeout int, comms Comms) RaftFSM {
	fsm := RealRaftFSM{
		id:          id,
		role:        Follower,
		currentTick: 0,
		peers:       peers,
		log:         NewInMemoryLog(),
		logger:      log.New(os.Stdout, fmt.Sprintf("[raft-%d]", id), log.Ltime),
		comms:       comms,

		CommonState: CommonState{
			VotedFor:         -1,
			CommitIndex:      0,
			LastAppliedIndex: 0,
			CurrentTerm:      0,
		},

		CandidateState: CandidateState{
			VotesReceived:    make(map[ID]bool),
			ElectionTimeout:  electionTimeout,
			ElectionDeadline: electionTimeout,
			Majority:         (len(peers) / 2) + 1,
		},

		FollowerState: FollowerState{
			HeartBeatTimeout:  heartBeatTimeout,
			HeartBeatDeadline: heartBeatTimeout,
		},

		LeaderState: LeaderState{
			NextIndices:  make(map[ID]int),
			MatchIndices: make(map[ID]int),
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

func (r *RealRaftFSM) GetComms() Comms {
	return r.comms
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
			r.logger.Println("heartbeat timeout.")
			r.BecomeCandidate()
			return
		}
	case Candidate:
		if r.currentTick >= r.CandidateState.ElectionDeadline {
			r.logger.Println("election timeout, no leader.")
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
	r.logger.Println("becoming follower")
	r.role = Follower
	r.FollowerState.HeartBeatDeadline = r.currentTick + r.FollowerState.HeartBeatTimeout
}

func (r *RealRaftFSM) BecomeCandidate() {
	r.logger.Println("becoming candidate")
	r.CandidateState.ElectionDeadline = r.currentTick + r.ElectionTimeout
	r.role = Candidate
	r.CurrentTerm++
	r.ResetVotes()
	r.GotVoteFrom(r.id)

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
	r.logger.Println("becoming leader")
	r.role = Leader
	r.ResetNextHeartBeatTime(r.currentTick)

	term, index := r.log.GetLastLogTermIndex()

	// send initial heart beat to all other nodes.
	msg := rpc.AppendEntriesReq{
		Term:              r.CommonState.CurrentTerm,
		LeaderId:          int(r.id),
		LeaderCommitIndex: r.CommitIndex(),
		PrevLogIndex:      index,
		PrevLogTerm:       term,
		Entries:           nil,
	}

	r.broadcastRpcImmediate(&msg)
}

func (r *RealRaftFSM) ReceiveMsg(msg rpc.Message) {
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
	r.logger.Println("received vote request from", msg.CandidateId)
	_, logIndex := r.log.GetLastLogTermIndex()
	if r.id == ID(msg.CandidateId) {
		return
	}

	voteReply := &rpc.VotedFor{
		VoterId:     int(r.id),
		Term:        r.CurrentTerm,
		VoteGranted: false,
	}

	r.logger.Println("sending vote for reply")

	if msg.Term < r.CurrentTerm {
		voteReply.VoteGranted = false
		r.sendRpcImmediate(ID(msg.CandidateId), voteReply)
		return
	}

	if r.VotedFor() < 0 {
		voteReply.VoteGranted = true
		r.sendRpcImmediate(ID(msg.CandidateId), voteReply)
		return
	}

	if r.VotedFor() == ID(msg.CandidateId) {
		voteReply.VoteGranted = msg.LastLogIndex >= logIndex
		r.sendRpcImmediate(ID(msg.CandidateId), voteReply)
		return
	}

	r.sendRpcImmediate(ID(msg.CandidateId), voteReply)
}

func (r *RealRaftFSM) ReceiveAppendEntriesRpc(msg *rpc.AppendEntriesReq) {
	senderId := msg.LeaderId
	_, logIndex := r.log.GetLastLogTermIndex()

	if ID(senderId) == r.id {
		return
	}

	if msg.Term < r.CurrentTerm {
		reply := rpc.AppendEntriesReply{
			Success: false,
			Term:    r.CurrentTerm,
		}
		r.sendRpcImmediate(ID(msg.LeaderId), &reply)
		return
	}

	if msg.Term > r.CurrentTerm {
		r.CommonState.CurrentTerm = msg.Term
		if r.role != Follower {
			r.BecomeFollower()
		}
	}

	if msg.LeaderCommitIndex > r.CommitIndex() {
		if msg.LeaderCommitIndex < logIndex {
			r.CommonState.CommitIndex = msg.LeaderCommitIndex
		} else {
			r.CommonState.CommitIndex = logIndex
		}
	}

	r.ResetHeartBeatTimeout(r.currentTick)

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

	reply := rpc.AppendEntriesReply{
		Term:    r.CurrentTerm,
		Success: true,
	}

	r.sendRpcImmediate(ID(senderId), &reply)
}
func (r *RealRaftFSM) ReceiveVotedFor(msg *rpc.VotedFor) {
	r.logger.Println("received vote from", msg.VoterId)
	if r.role != Candidate {
		r.unhandledRpc(msg)
		return
	}

	if !msg.VoteGranted {
		return
	}

	r.GotVoteFrom(ID(msg.VoterId))

	if r.CandidateState.HasMajority() {
		r.BecomeLeader()
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
	if r.LeaderState.NextHeartBeatTime < r.currentTick {

		msg := &rpc.AppendEntriesReq{
			Term:              r.CommonState.CurrentTerm,
			LeaderId:          int(r.id),
			LeaderCommitIndex: r.CommonState.CommitIndex,
			PrevLogIndex:      r.CommonState.LastAppliedIndex,
			Entries:           nil,
		}
		r.broadcastRpcImmediate(msg)

		r.ResetNextHeartBeatTime(r.currentTick)
		return
	}

}

func (r *RealRaftFSM) followerStep() {
	// 1. check heartbeat timeout
}

func (r *RealRaftFSM) enqueueBroadcastRpc(msg rpc.Message) {
	go r.comms.BroadcastRpc(context.Background(), msg)
}

func (r *RealRaftFSM) broadcastRpcImmediate(msg rpc.Message) {
	r.logger.Println("broadcasting rpc", msg, reflect.TypeOf(msg))
	r.comms.BroadcastRpc(context.Background(), msg)
}

func (r *RealRaftFSM) sendRpcImmediate(receiver ID, msg rpc.Message) {
	r.logger.Println("sending rpc to", receiver)
	r.comms.Rpc(context.Background(), receiver, msg)
}

func (r *RealRaftFSM) enqueueRpc(receiver ID, msg rpc.Message) {
	go r.comms.Rpc(context.Background(), receiver, msg)
}

func (r *RealRaftFSM) ReplicateToLog(data interface{}) (error) {
	prevLogTerm, prevLogIndex := r.log.GetLastLogTermIndex()

	entry := Entry{
		Index: -1,
		Term:  r.CurrentTerm,
		Data:  data,
	}
	err := r.log.Append(entry)
	if err != nil {
		return err
	}

	msg := rpc.AppendEntriesReq{
		LeaderId:          int(r.id),
		LeaderCommitIndex: r.CommitIndex(),

		Term:         r.CurrentTerm,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      []interface{}{data},
	}

	r.broadcastRpcImmediate(&msg)

	return err
}
