package kvserver

import (
	"fmt"
	"github.com/zl14917/MastersProject/pkg/kvstore"
	"github.com/zl14917/MastersProject/pkg/logger"
	"github.com/zl14917/MastersProject/router"
	"google.golang.org/grpc"
	"log"
	"os"
	"path"
	"strconv"
)

const lockFileName = "cliftondb.lock.file"
const logPath = "logs/"
const metaPath = "metadata/"
const dataPath = "data/"
const PartionLockFileName = "partition.lock.file"

const logPrefix = "[kv-server]"
const logFileName = "kvserver.log"

type PartitionId int
type KVServer struct {
	Conf Config

	kvGrpcApiServer *GrpcKVServer

	kvStores map[PartitionId]*kvstore.CliftonDBKVStore

	requestRouter *router.ClientRequestRouter
	grpcServer    grpc.Server

	DbPath       string
	LockFilePath string
	LogsPath     string
	MetadatPath  string
	DataPath     string

	Logger *logger.FileLogger
}

func NewKVServer(conf Config) (*KVServer, error) {
	dbPath := conf.DbPath
	logsPath := path.Join(dbPath, logPath, logFileName)

	fileLogger, err := logger.NewFileLogger(path.Join(dbPath, logPath), logPath, log.LstdFlags)
	if err != nil {
		return nil, err
	}

	return &KVServer{
		Conf: conf,

		kvGrpcApiServer: nil,
		kvStores:        make(map[PartitionId]*kvstore.CliftonDBKVStore),
		requestRouter:   nil,

		DbPath:       dbPath,
		LockFilePath: path.Join(dbPath, lockFileName),
		LogsPath:     logsPath,
		MetadatPath:  path.Join(dbPath, metaPath),
		DataPath:     path.Join(dbPath, dataPath),

		Logger: fileLogger,
	}, nil
}

func (s *KVServer) BootstrapServer() error {
	if len(s.Conf.Nodes.PeerList) < 1 {
		return s.boostrapInStandaloneMode()
	}
	return s.boostrapInClusterMode()
}

func (s *KVServer) detectLockFile() (bool, error) {
	var (
		stat os.FileInfo
		err  error
	)
	if stat, err = os.Stat(s.LockFilePath); err != nil {
		return false, err
	}

	return !stat.IsDir(), nil
}

func (s *KVServer) writeLockFile() error {
	file, err := os.Create(s.LockFilePath)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (s *KVServer) boostrapInClusterMode() error {
	return nil
}

// Standalone mode startup procedures:
// 1. check if lock file exists.
// 2. if so, proceed,
// 3. if not create all folders
// load or create kvstore partitions from Conf.
func (s *KVServer) boostrapInStandaloneMode() error {

	//
	var (
		prevInstance bool
		err          error
		partitionIds []PartitionId
	)
	prevInstance, err = s.detectLockFile()

	if !prevInstance {
		const DefaultPartitionId = 0
		err = s.createFolders()
		if err != nil {
			return fmt.Errorf("cannot create all folders:", err)
		}
		partitionIds = []PartitionId{DefaultPartitionId}

	} else {
		partitionIds, err = getPartitionsFromDataDir(s.DataPath, PartionLockFileName)
		if err != nil {
			return err
		}
	}

	err = s.boostrapKVStoresForEachPartition(partitionIds)

	if err != nil {
		return err
	}

	if !prevInstance {
		err = s.writeLockFile()
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *KVServer) boostrapKVStoresForEachPartition(partitionIds []PartitionId) error {
	for _, id := range partitionIds {
		storeRootPath := path.Join(s.DataPath, strconv.Itoa(int(id)))
		store, err := kvstore.NewCliftonDBKVStore(storeRootPath)

		if err != nil {
			return err
		}

		s.kvStores[id] = store
	}
	return nil
}

func (s *KVServer) createFolders() error {
	return ensureDirsExist(s.DataPath, s.MetadatPath, s.LogsPath)
}

func (s *KVServer) LoadPartitions() error {

}
