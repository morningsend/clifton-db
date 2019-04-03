package cliftondbserver

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
)

type ReadConsistencyLevel int

const (
	ReadLocalCopy ReadConsistencyLevel = iota
)

const (
	DefaultReadConsistency = ReadLocalCopy
)

type KvStorePut struct {
	KeyHash uint32
	Key     string
	Value   []byte
}

type KvStoreGet struct {
	Key string
}

type KvStoreDelete struct {
	KeyHash uint32
	Key     string
}

type ReplicatedKvStore struct {
	raftNode *RaftNode
	proposeC chan<- string
	commitC  <-chan *string
}

func (p *ReplicatedKvStore) NewReplicatedKvStore() (*ReplicatedKvStore, error) {

	return nil, errors.New("not implemented")
}

func (b *ReplicatedKvStore) ProposePut(ctx context.Context, key string, value []byte) error {
	var buf bytes.Buffer
	putReq := KvStorePut{
		Key:   key,
		Value: value,
	}
	if err := gob.NewEncoder(&buf).Encode(&putReq); err != nil {
		return err
	}

	b.proposeC <- buf.String()
	return nil
}

func (st *ReplicatedKvStore) ProposeDelete(ctx context.Context, key string) error {
	return nil
}

func (st *ReplicatedKvStore) ProposeGet(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
