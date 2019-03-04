package sstable

type MergeCallback func(success bool, err error)

func MergeTables(source []SSTable, destination *SSTable, callback MergeCallback) {

}
