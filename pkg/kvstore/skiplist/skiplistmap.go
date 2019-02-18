package skiplist

type Key string
type Value []byte

type MapIterator interface {
	Next() bool
	Current() (key Key, value Value)
}

type Map interface {
	Get(key Key) (value Value, ok bool)
	Put(key Key, value Value) (err error)
	Contains(key Key) (ok bool)
	Iterator() MapIterator
}

type SkipListMap struct {
	headSentinel *slMapNode
	head         *slMapNode

	size     int
	maxLevel int
}

type slMapNode struct {
	Key   Key
	Value Value
	Next  []*slMapNode
}

type SkipListIterator struct {
	current *slMapNode
}

func makeMapNode(key Key, value Value, height int) *slMapNode {
	return &slMapNode{
		Key:   key,
		Value: value,
		Next:  make([]*slMapNode, 1, height),
	}
}

func NewSkipListMap() Map {
	maxLevel := 6
	sentinel := makeMapNode("", nil, maxLevel)
	return &SkipListMap{
		size:         0,
		headSentinel: sentinel,
		head:         sentinel,
		maxLevel:     maxLevel,
	}
}

func (m *SkipListMap) Get(key Key) (value Value, ok bool) {
	node, exists := m.searchNode(key)
	if ! exists {
		return nil, false
	}
	return node.Value, true
}

func (m *SkipListMap) Contains(key Key) (bool) {
	_, exists := m.searchNode(key)
	return exists
}

func (m *SkipListMap) searchNode(key Key) (node *slMapNode, exists bool) {
	return
}

func (m *SkipListMap) linkNode(previous *slMapNode, newNode *slMapNode) {

}
func (m *SkipListMap) Put(key Key, value Value) (err error) {
	if m.head == nil {
		m.head = makeMapNode(key, value, m.maxLevel)
		return
	}
	node, exists := m.searchNode(key)
	if exists {
		node.Value = value
		return
	}

	previous := node
	newNode := makeMapNode(key, value, m.maxLevel)
	m.linkNode(previous, newNode)
	return nil
}

func (m *SkipListMap) Iterator() MapIterator {
	return &SkipListIterator{
		current: m.head,
	}
}

func (i *SkipListIterator) Next() bool {
	if len(i.current.Next) > 0 && i.current.Next[0] != nil {
		i.current = i.current.Next[0]
		return true
	}
	return false
}

func (i *SkipListIterator) Current() (key Key, value Value) {
	if i.current == nil {
		return
	}
	return i.current.Key, i.current.Value
}
