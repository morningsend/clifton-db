package kvserver

import (
	"github.com/zl14917/MastersProject/kvstore"
	"github.com/zl14917/MastersProject/logger"
	"github.com/zl14917/MastersProject/router"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
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

type KvServerLockFileData struct {
	Partitions []PartitionId
}

var defaultKvServerLockFileData = KvServerLockFileData{
	Partitions: []PartitionId{0},
}

type KVServer struct {
	Conf Config

	kvGrpcApiServer *GrpcKVServer

	Partitions []PartitionId
	kvStores   map[PartitionId]*kvstore.CliftonDBKVStore

	requestRouter *router.ClientRequestRouter
	grpcServer    grpc.Server

	DbPath       string
	LockFilePath string
	LogsPath     string
	MetadatPath  string
	DataPath     string

	Logger logger.Logger
}

func (s *KVServer) RegisterGrpcServer(server *grpc.Server) {
	s.kvGrpcApiServer.Register(server)
}

func NewKVServer(conf Config) (*KVServer, error) {
	var serverLogger logger.Logger
	dbPath := conf.DbPath
	logsPath := path.Join(dbPath, logPath, logFileName)

	serverLogger, err := logger.NewFileLogger(path.Join(dbPath, logPath), logPrefix, log.LstdFlags)
	if err != nil {
		serverLogger = log.New(os.Stdout, logPrefix, log.LstdFlags)
	}

	server := &KVServer{
		Conf: conf,

		kvGrpcApiServer: nil,
		kvStores:        make(map[PartitionId]*kvstore.CliftonDBKVStore),
		requestRouter:   nil,

		DbPath:       dbPath,
		LockFilePath: path.Join(dbPath, lockFileName),
		LogsPath:     logsPath,
		MetadatPath:  path.Join(dbPath, metaPath),
		DataPath:     path.Join(dbPath, dataPath),

		Logger: serverLogger,
	}

	return server, nil
}

func (s *KVServer) BootstrapServer() error {
	if len(s.Conf.Nodes.PeerList) < 1 {
		return s.boostrapInStandaloneMode()
	}
	return s.boostrapInClusterMode()
}

func (s *KVServer) detectLockFile() (KvServerLockFileData, error) {
	var (
		err  error
		data KvServerLockFileData
	)
	file, err := os.OpenFile(s.LockFilePath, os.O_RDONLY, 0644)

	if err != nil {
		return defaultKvServerLockFileData, nil
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&data)
	if err != nil {
		return defaultKvServerLockFileData, err
	}
	return data, nil
}

func (s *KVServer) WriteLockFile(data KvServerLockFileData) error {
	file, err := os.OpenFile(s.LockFilePath, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	err = yaml.NewEncoder(file).Encode(data)

	if err != nil {
		file.Close()
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

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
		err error
	)
	savedSettings, err := s.detectLockFile()

	s.Partitions = make([]PartitionId, len(savedSettings.Partitions))

	copy(s.Partitions, savedSettings.Partitions)

	err = s.boostrapKVStoresForEachPartition(savedSettings.Partitions)

	if err != nil {
		return err
	}

	err = s.WriteLockFile(savedSettings)

	if err != nil {
		return err
	}

	return nil
}

func (s *KVServer) boostrapKVStoresForEachPartition(partitionIds []PartitionId) error {
	for _, id := range partitionIds {
		storeRootPath := path.Join(s.DataPath, strconv.Itoa(int(id)))
		store, err := kvstore.NewCliftonDBKVStore(storeRootPath, s.LogsPath)

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
