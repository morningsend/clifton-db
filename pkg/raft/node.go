package raft

type RaftNode struct {
	Id  int
	FSM RaftFSM
}
