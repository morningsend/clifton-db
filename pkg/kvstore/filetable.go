package kvstore

type FileTable interface {
	BeginFlushing()
}
