package cliftondbserver

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/api/cluster-services"
	"github.com/zl14917/MastersProject/kvstore"
	"github.com/zl14917/MastersProject/router"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"path"
	"strconv"
)

const lockFileName = "cliftondb.lock.file"
const logPath = "logs/"
const metaPath = "metadata/"
const partitionPath = "partitions/"
const PartionLockFileName = "partition.lock.file"

const logPrefix = "[kv-grpcServer]"
const logFileName = "kvserver.log"

type PartitionId int

type KvServerLockFileData struct {
	Partitions []PartitionId
}

var defaultKvServerLockFileData = KvServerLockFileData{
	Partitions: []PartitionId{0},
}

type CliftonDbServer struct {
	Conf Config

	Partitions []PartitionId
	kvStores   map[PartitionId]*kvstore.CliftonDBKVStore

	grpcServer    *grpc.Server
	requestRouter *router.ClientRequestRouter

	DbRootPath    string
	LockFilePath  string
	LogsPath      string
	MetadatPath   string
	PartitionPath string
	Logger        *zap.Logger

	listener *StoppableListener
}

func (s *CliftonDbServer) AddNode(context.Context, *cluster_services.AddNodeReq) (*cluster_services.AddNodeRes, error) {
	panic("implement me")
}

func (s *CliftonDbServer) GetNodeList(context.Context, *interface{}) (*cluster_services.PeerListRes, error) {
	panic("implement me")
}

func (s *CliftonDbServer) RemoveNode(context.Context, *cluster_services.RemoveNodeReq) (*cluster_services.RemoveNodeRes, error) {
	panic("implement me")
}

func NewCliftonDbServer(conf Config) (*CliftonDbServer, error) {
	dbPath := conf.DbPath
	logsPath := path.Join(dbPath, logPath, logFileName)

	serverLogger, err := zap.
		NewDevelopmentConfig().
		Build()

	if err != nil {
		serverLogger = zap.NewExample()
	}

	server := &CliftonDbServer{
		Conf: conf,

		kvStores:      make(map[PartitionId]*kvstore.CliftonDBKVStore),
		requestRouter: nil,

		DbRootPath:    dbPath,
		LockFilePath:  path.Join(dbPath, lockFileName),
		LogsPath:      logsPath,
		MetadatPath:   path.Join(dbPath, metaPath),
		PartitionPath: path.Join(dbPath, partitionPath),

		Logger: serverLogger,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Conf.Server.ListenPort))

	if err != nil {
		return nil, err
	}

	sl, err := NewStoppableListener(listener)
	server.listener = sl

	if err != nil {
		serverLogger.Error("error starting listener", zap.Error(err))
		err = listener.Close()
		return nil, err
	}

	server.ServeClusterNodeApi()

	return server, nil
}

func (s *CliftonDbServer) Boostrap() error {
	if len(s.Conf.Nodes.PeerList) < 1 {
		return s.boostrapInStandaloneMode()
	}
	return s.boostrapInClusterMode()
}

func (s *CliftonDbServer) detectLockFile() (KvServerLockFileData, error) {
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

func (s *CliftonDbServer) WriteLockFile(data KvServerLockFileData) error {
	file, err := os.OpenFile(s.LockFilePath, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	err = yaml.NewEncoder(file).Encode(data)

	if err != nil {
		ferr := file.Close()
		if ferr != nil {
			s.Logger.Error("error closing file", zap.Error(ferr))
		}
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *CliftonDbServer) boostrapInClusterMode() error {
	s.Logger.Info("bootstraping cliftondb grpcServer in cluster mode")
	return nil
}

// Standalone mode startup procedures:
// 1. check if lock file exists.
// 2. if so, proceed,
// 3. if not create all folders
// load or create kvstore partitions from Conf.
func (s *CliftonDbServer) boostrapInStandaloneMode() error {
	s.Logger.Info("boostraping cliftondb grpcServer in standalone mode")
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

func (s *CliftonDbServer) boostrapKVStoresForEachPartition(partitionIds []PartitionId) error {
	for _, id := range partitionIds {
		storeDirPath := path.Join(s.PartitionPath, strconv.Itoa(int(id)))
		s.Logger.Info("will create/open kv-store at path", zap.String("storeDirPath", storeDirPath))
		store, err := kvstore.NewCliftonDBKVStore(storeDirPath, s.LogsPath)

		if err != nil {
			return err
		}
		s.kvStores[id] = store
	}
	return nil
}

func (s *CliftonDbServer) createFolders() error {
	return ensureDirsExist(s.PartitionPath, s.MetadatPath, s.LogsPath)
}

func (s *CliftonDbServer) LookupPartitions(key string) (kv *kvstore.CliftonDBKVStore, ok bool) {
	return nil, false
}

func (s *CliftonDbServer) ServeClusterNodeApi() {
	cluster_services.RegisterClusterNodeServer(s.grpcServer, s)
}

func (s *CliftonDbServer) Shutdown() {
	if s.listener != nil {
		s.listener.Stop()
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
