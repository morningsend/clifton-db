package server

import "path"

type RaftOptions struct {
	RaftWALDirPath   string
	RaftSnapshotPath string
	RaftNodeID       int64
}

type LocalStorageOptions struct {
	LocalStorageDirPath string
	Engine              string
	BadgerEngineOptions string
}

type Options struct {
	RaftOptions
	LocalStorageOptions
}

func DefaultServerOptions(raftId int64, localDbPath string) Options {
	return Options{
		RaftOptions: RaftOptions{
			RaftWALDirPath:   path.Join(localDbPath, "raft/wal"),
			RaftSnapshotPath: path.Join(localDbPath, "raft/snapshot"),
		},
		LocalStorageOptions: LocalStorageOptions{
			LocalStorageDirPath: path.Join(localDbPath, "db"),
			Engine:              "badger",
		},
	}
}
