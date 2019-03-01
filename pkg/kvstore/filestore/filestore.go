package filestore

import "github.com/zl14917/MastersProject/pkg/kvstore/sstable"

type FileStore interface {
}

type LevelFileStore struct {
	SizeBytes uint32
	Level0Tables []sstable.SSTable
}

func NewLevelFileStore() FileStore {
	return nil
}
