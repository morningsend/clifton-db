package tables

type FileTable interface {
	BeginFlushing(table MemTable) error
}
