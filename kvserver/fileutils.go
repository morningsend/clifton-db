package kvserver

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func ensureDirsExist(folderPaths ...string) error {
	var err error
	for _, p := range folderPaths {
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPartitionsFromDataDir(dataDirPath string, lockFileName string) (partIds []PartitionId, err error) {
	files, err := ioutil.ReadDir(dataDirPath)
	if err != nil {
		return nil, err
	}

	partIds = make([]PartitionId, 0, len(files))

	for _, file := range files {
		if ok, _ := IsFileExists(path.Join(dataDirPath, file.Name(), lockFileName)); !ok {
			continue
		}

		id, err := strconv.Atoi(file.Name())

		if err != nil {
			continue
		}
		partIds = append(partIds, PartitionId(id))
	}

	return
}

func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return !os.IsNotExist(err), err
	}
	return true, nil
}
