package dht

import "math"

const NoPartition = -1

var (
	wholeSegment HashSegments = HashSegments{
		RangeStart:  math.MinInt32,
		RangeEnd:    math.MaxInt32,
		PartitionId: NoPartition,
	}
)

type HashSegments struct {
	RangeStart  int32
	RangeEnd    int32
	PartitionId int32
}

type PartitionTable struct {
	PartitionId int
	Nodes       []NodeId
}

type NodeId struct {
	Id int
}

// Hash ring maps intervals of segments into logical partitions.
// It supports operation to look up a node in O(log(n)) partitions

type HashRing struct {
	root *HashSegments
}

func NewHashRing() HashRing {
	rootSegment := wholeSegment
	return HashRing{
		root: &rootSegment,
	}
}
