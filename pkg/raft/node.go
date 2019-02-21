package raft

type RaftNode struct {
	Id     int
	FSM    RaftFSM
	Comms  Comms
	Driver RaftDriver
}


func NewNode() *RaftNode {
	return nil
}

func Start() {

}
