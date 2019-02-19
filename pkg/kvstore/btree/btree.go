package btree

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/blockstore"
	"github.com/zl14917/MastersProject/pkg/kvstore/maps"
	"sort"
	"strings"
)

// B+ Tree implementation
// Currently only support memory storage
type BTreeMap struct {
	maxFanout     int
	minFanout     int
	maxLeafFanout int
	minKeys       int
	root          *internalNode
	storage       blockstore.BlockStorage
}

type internalNode struct {
	Keys          []string
	parent        *internalNode
	InternalNodes []*internalNode
	LeafNodes     []*leafNode
}

type BlockRecord struct {
	Removed bool
	Data    []byte
}

type leafNode struct {
	Parent   *internalNode
	Pointers map[maps.Key]blockstore.BlockPointer
	Next     *leafNode
}

type block struct {
	data []byte
}

type fileHeader struct {
	blocks    int
	blockSize int
}

func NewBTreeMap(fanoutFactor int, blockSize int) maps.Map {
	root := &internalNode{
		Keys:   make([]string, 0, fanoutFactor),
		parent: nil,

		InternalNodes: nil,
		LeafNodes:     make([]*leafNode, 0, fanoutFactor),
	}

	leaf := &leafNode{
		Pointers: make(map[maps.Key]blockstore.BlockPointer),
		Next:     nil,
		Parent:   root,
	}

	root.LeafNodes = append(root.LeafNodes, leaf)

	return &BTreeMap{
		minFanout:     fanoutFactor,
		maxFanout:     fanoutFactor * 2,
		maxLeafFanout: fanoutFactor * 2,
		root:          root,

		storage: blockstore.NewInMemoryStore(blockSize),
	}
}

func FromBytes(data []byte) maps.Map {
	return nil
}

func (t *BTreeMap) Get(key maps.Key) (value maps.Value, ok bool) {
	leaf := t.search(key)

	//var dp BlockPointer

	dp, found := leaf.Pointers[key]

	if !found {
		return nil, false
	}

	data, err := t.storage.ReadBytes(dp)

	if err != nil {
		return nil, false
	}

	return data, true
}

func (t *BTreeMap) splitLeaf(node *leafNode) {
	if len(node.Pointers) < t.maxLeafFanout {
		return
	}
	keys := make([]string, 0, len(node.Pointers))
	for k, _ := range node.Pointers {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	mid := len(node.Pointers) / 2

	leftKeys := keys[:mid]
	rightKeys := keys[mid:]

	leftLeaf := node
	rightLeaf := &leafNode{
		Pointers: make(map[maps.Key]blockstore.BlockPointer),
		Next:     leftLeaf.Next,
		Parent:   leftLeaf.Parent,
	}

	leftLeaf.Next = rightLeaf

	for _, k := range rightKeys {
		rightLeaf.Pointers[maps.Key(k)] = leftLeaf.Pointers[maps.Key(k)]
		delete(leftLeaf.Pointers, maps.Key(k))
	}
	leftKey := leftKeys[len(leftKeys)-1]
	// rightKey = old key in parent corresponding to node.
	rightKey := rightKeys[len(rightKeys)-1]
	p := leftLeaf.Parent

	p.LeafNodes = append(p.LeafNodes, nil)
	p.Keys = append(p.Keys, "")

	var oldPos int
	for i := len(p.Keys) - 1; i > 0; i-- {
		p.LeafNodes[i+1] = p.LeafNodes[i]
		p.Keys[i] = p.Keys[i-1]

		k := p.Keys[i]
		oldPos = i
		if strings.Compare(string(k), rightKey) <= 0 {
			break
		}
	}
	p.LeafNodes[oldPos] = rightLeaf
	p.LeafNodes[oldPos-1] = leftLeaf
	p.Keys[oldPos] = rightKey
	p.Keys[oldPos-1] = leftKey
}

func (t *BTreeMap) splitInternal(node *internalNode) {

}
func (t *BTreeMap) Contains(key maps.Key) (ok bool) {
	return
}

func (t *BTreeMap) Put(key maps.Key, value maps.Value) (err error) {
	leaf := t.search(key)

	bp, err := t.storage.WriteBytes(value, len(value))
	if err != nil {
		return err
	}

	leaf.Pointers[key] = bp
	if len(leaf.Pointers) > t.maxFanout {
		t.splitLeaf(leaf)
	}
	return nil
}

func (t *BTreeMap) Iterator() maps.MapIterator {
	return nil
}

func (t *BTreeMap) search(key maps.Key) (leaf *leafNode) {
	var (
		current *internalNode = t.root
	)

	var keyIndex int

	for current.LeafNodes == nil {
		keyIndex = len(current.InternalNodes)
		for i, k := range current.Keys {
			if strings.Compare(string(k), string(key)) <= 0 {
				keyIndex = i
			}
		}
		current = current.InternalNodes[keyIndex]
	}

	keyIndex = len(current.Keys)

	for i, k := range current.Keys {
		if strings.Compare(string(k), string(key)) <= 0 {
			keyIndex = i
		}
	}

	return current.LeafNodes[keyIndex]
}
