package raft

import (
	"github.com/zl14917/MastersProject/pkg/raft/rpc"
	"log"
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
	Logs        []interface{}

	CommitIndex      int
	LastAppliedIndex int
}

type CandidateState struct {
	Majority         int
	ElectionDeadline int
	VotesReceived    map[ID]bool
}

type LeaderState struct {
}

type FollowerState struct {
	HeartBeatDeadline int
}

type RaftFSM interface {
	raftFSM()
	Id() ID
	// states
	Role() Role
	CurrentTerm() int
	CommitIndex() int
	LastAppliedIndex() int
	VotedFor() ID
	Logs() []interface{}
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
	logger           log.Logger
	peers            []ID
}

func NewRaftFSM(id int, peers []ID, electionTimeout int, heartBeatTimeout int) RaftFSM {
	fsm := RealRaftFSM{
		id:               ID(id),
		role:             Follower,
		currentTick:      0,
		electionTimeout:  electionTimeout,
		heartBeatTimeout: heartBeatTimeout,
		peers:            peers,
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

func (r *RealRaftFSM) CurrentTerm() int {
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

func (r *RealRaftFSM) Logs() []interface{} {
	return nil
}

func (r *RealRaftFSM) tick() (int) {
	r.currentTick++
	r.checkTimeout()
	r.step()
	return r.currentTick
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
	r.VotesReceived = make(map[ID]bool, len(r.peers))

	requestVoteMsg := &rpc.RequestVote{
		CandidateId: int(r.id), // TODO: refactor ID type to its own package,
		Term: r.CurrentTerm(),
		LastLogIndex: 0,
		LastLogTerm: 0,
	}
	r.enqueueBroadcastRpc(requestVoteMsg)
}

func (r *RealRaftFSM) BecomeLeader() {

}

func (r *RealRaftFSM) ReceiveRpc(msg rpc.Message) {
	switch x := msg.(type) {
	case *rpc.AppendEntriesReq:
		if r.role != Follower {
			r.unhandledRpc(r.role, x)
			return
		}
		r.ReceiveAppendEntriesRpc(x)
	case *rpc.AppendEntriesReply:
		if r.role != Leader {
			r.unhandledRpc(r.role, x)
			return
		}
		r.ReceiveAppendEntriesReply(x)
	case *rpc.RequestVote:
		if r.role != Follower {
			r.unhandledRpc(r.role, x)
			return
		}
		r.ReceiveVoteRequestRpc(x)
	case *rpc.VotedFor:
		if r.role != Candidate {
			r.unhandledRpc(r.role, x)
			return
		}

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

func (r *RealRaftFSM) ReceiveAppendEntriesRpc(msg *rpc.AppendEntriesReq) {

}
func (r *RealRaftFSM) ReceiveVotedFor(msg *rpc.VotedFor) {

}

func (r *RealRaftFSM) unhandledRpc(role Role, msg rpc.Message) {

}

func (r *RealRaftFSM) candidateStep() {
	votes := len(r.CandidateState.VotesReceived)
	if votes >= r.CandidateState.Majority {
		r.BecomeLeader()
	}
}

func (r *RealRaftFSM) leaderStep() {

}

func (r *RealRaftFSM) followerStep() {
	// 1. check heartbeat timeout

}

func (r *RealRaftFSM) enqueueBroadcastRpc(msg rpc.Message) {

}

func (r *RealRaftFSM) broadcastRpcImmediate(msg rpc.Message) {

}

func (r* RealRaftFSM) sendRpcImmediate(msg rpc.Message) {

}

func (r *RealRaftFSM) enqueueRpc(receiver ID, msg rpc.Message) {

}
