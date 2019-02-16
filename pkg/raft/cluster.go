package raft

type Cluster struct {
	AllMembers map[ID]bool
	SelfID     ID
}

func NewCluster(self ID, allIds []ID) *Cluster {
	members := make(map[ID]bool)
	for _, id := range allIds {
		members[id] = true
	}
	return &Cluster{
		SelfID:     self,
		AllMembers: members,
	}
}

func (c *Cluster) AddNode(id ID) {
	c.AllMembers[id] = true
}

func (c *Cluster) RemoveNode(id ID) {
	delete(c.AllMembers, id)
}
