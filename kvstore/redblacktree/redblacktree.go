package redblacktree

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/hash"
	"strings"
)

type rbColor int32

const (
	Black rbColor = iota
	Red
)

type Tree interface {
	Insert(key string, data []byte) error
	Remove(key string) (data []byte, err error)
	Size() int
	Exists(key string) (ok bool, err error)
}

type Iterator interface {
	Next() bool
	KVPair(key string, data []byte)
}

type rbTree struct {
	root   *rbTreeNode
	hasher hash.Hasher
	size   int
}

type rbTreeNode struct {
	color   rbColor
	key     string
	keyHash int32
	data    []byte

	left  *rbTreeNode
	right *rbTreeNode
}

func NewTree() Tree {
	return &rbTree{
		root:   nil,
		size:   0,
		hasher: hash.NewSimpleKeyHasher(),
	}
}

func (t *rbTree) Insert(key string, data []byte) error {
	if t.root == nil {
		t.root = t.newRBTreeNode(key, data, Black)
		return nil
	}

	var (
		x   = t.root
		y   *rbTreeNode
		cmp int = 0
	)

	for x != nil {
		y = x
		cmp = strings.Compare(key, x.key)
		if cmp == 0 {
			x.data = data
			return nil
		} else if cmp < 0 {
			x = x.left
		} else {
			x = x.right
		}
	}

	newNode := t.newRBTreeNode(key, data, Red)

	if cmp < 0 {
		y.left = newNode
	} else {
		y.right = newNode
	}
	y.fixUpChildren(nil)
	return nil
}

func (t *rbTree) Remove(key string) ([]byte, error) {
	return nil, nil
}

func (t *rbTree) Exists(key string) (ok bool, err error) {
	return
}

func (t *rbTree) Size() int {
	return 0
}

func (tree *rbTree) newRBTreeNode(key string, data []byte, color rbColor) *rbTreeNode {
	tree.hasher.UpdateString(key)
	keyHash := tree.hasher.GetHashInt32()

	return &rbTreeNode{
		color:   color,
		key:     key,
		data:    data,
		keyHash: keyHash,
		left:    nil,
		right:   nil,
	}
}

func (node *rbTreeNode) leftRotate() *rbTreeNode {
	x := node
	y := node.right

	if y == nil {
		return node
	}

	x.right = y.left
	y.left = x

	return y
}

func (node *rbTreeNode) rightRotate() *rbTreeNode {
	y := node
	x := node.left

	if x == nil {
		return y
	}

	y.left = x.right
	x.right = y

	return x
}

func (node *rbTreeNode) fixUpChildren(parent *rbTreeNode) {

}

func (tree *rbTree) iterator() Iterator {

}
