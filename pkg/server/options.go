package server

import "path"

type RaftOptions struct {
	RaftWALDirPath   string
	RaftSnapshotPath string
	RaftNodeID       uint64
}

type LocalStorageOptions struct {
	LocalStorageDirPath string
}

type Options struct {
	RaftOptions
	LocalStorageOptions
}

func DefaultServerOptions(raftId uint64, localDbPath string) Options {
	return Options{
		RaftOptions: RaftOptions{
			RaftNodeID:       raftId,
			RaftWALDirPath:   path.Join(localDbPath, "raft/wal"),
			RaftSnapshotPath: path.Join(localDbPath, "raft/snapshot"),
		},
		LocalStorageOptions: LocalStorageOptions{
			LocalStorageDirPath: path.Join(localDbPath, "db"),
		},
	}
}
