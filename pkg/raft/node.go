package raft

type RaftNode struct {
	Id    int
	FSM   RaftFSM
	Comms Comms
	Driver Driver
}


func NewNode() *RaftNode {
	return nil
}

func Start() {

}
