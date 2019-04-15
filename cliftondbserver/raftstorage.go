package cliftondbserver

import (
	"go.etcd.io/etcd/etcdserver/api/snap"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/wal"
	"go.etcd.io/etcd/wal/walpb"
	"go.uber.org/zap"
	"io"
)

type RaftPersistentStorage interface {
	Save(state raftpb.HardState, entries [] raftpb.Entry) error
	SaveSnap(sp raftpb.Snapshot) error
	Close() error
}

type raftPersistentStorage struct {
	*wal.WAL
	*snap.Snapshotter
}

func NewRaftPersistentStorage(w *wal.WAL, s *snap.Snapshotter) *raftPersistentStorage {
	return &raftPersistentStorage{
		WAL:         w,
		Snapshotter: s,
	}
}

var _ RaftPersistentStorage = &raftPersistentStorage{}

func (s *raftPersistentStorage) SaveSnap(snapshot raftpb.Snapshot) error {
	walsnap := walpb.Snapshot{
		Index: snapshot.Metadata.Index,
		Term:  snapshot.Metadata.Term,
	}
	err := s.WAL.SaveSnapshot(walsnap)
	if err != nil {
		return err
	}

	err = s.Snapshotter.SaveSnap(snapshot)
	if err != nil {
		return nil
	}

	return s.WAL.ReleaseLockTo(snapshot.Metadata.Index)
}

func readWALContent(lg *zap.Logger, waldir string, sp walpb.Snapshot) (w *wal.WAL, s raftpb.HardState, entries []raftpb.Entry) {
	var
	(
		err           error
		//metadata      []byte
		triedToRepair bool
	)

	for {
		if w, err = wal.Open(lg, waldir, sp); err != nil {
			if lg != nil {
				lg.Fatal("cannot open wal", zap.String("wal-dir", waldir), zap.Error(err))
			}
		}

		_, s, entries, err = w.ReadAll()
		if err != nil {
			closeErr := w.Close()
			if closeErr != nil {
				lg.Error("error closing wal", zap.Error(err))
			}
			if triedToRepair || err != io.ErrUnexpectedEOF {
				lg.Fatal("cannot repair WAL, exiting", zap.Error(err))
			}

			if wal.Repair(lg, waldir) {
				triedToRepair = true
			}
			continue
		}
		break
	}
	return
}
