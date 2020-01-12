package engine

import "github.com/morningsend/clifton-db/pkg/storagepb"

type Timestamp = storagepb.Timestamp

type Engine interface {
	Reader()
	ReaderAt(timestamp Timestamp)
	Writer()
	WriterAt(timestamp Timestamp)
}
