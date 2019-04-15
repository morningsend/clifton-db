package cliftondbclient

import (
	"context"
	"errors"
	"github.com/zl14917/MastersProject/api/kv-client"
	"google.golang.org/grpc"
	"time"
)

var (
	EmptyServerListErr = errors.New("server list is empty")
	NotLeaderErr       = errors.New("not a leader")
	NotConnectedErr    = errors.New("not connected")
	PartitionUnavailable = errors.New("partition are not in quorum or leader maybe down")
)

type CliftonDbOps interface {
	Get(key string) (value []byte, err error)
	Put(key string, value [] []byte) (err error)
	Delete(key string) (err error)
}

type CliftonDbClient interface {
	Connect(serverList []string) error
	ConnectContext(ctx context.Context, serverList []string) error
	Close() error

	CliftonDbOps
}

type clientConf struct {
	clientId   uint64
	leaderAddr string
}

type cliftonDbClient struct {
	clientConf
	*grpc.ClientConn

	kvClient kv_client.ClientClient
	timeout  time.Duration
}

type CliftonDbOptions func(client *cliftonDbClient)

func New(options ...CliftonDbOptions) CliftonDbClient {
	c := &cliftonDbClient{}
	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *cliftonDbClient) ExecuteCmd(ctx context.Context, in *kv_client.ClientCmdReq, opts ...grpc.CallOption) (*kv_client.ClientCmdRes, error) {

}

func (c *cliftonDbClient) Query(ctx context.Context, in *kv_client.ClientQueryReq, opts ...grpc.CallOption) (*kv_client.ClientQueryRes, error) {
	panic("implement me")
}

func (c *cliftonDbClient) RegisterClient(ctx context.Context, in *kv_client.RegisterClientReq, opts ...grpc.CallOption) (*kv_client.RegisterClientRes, error) {
	panic("implement me")
}

func (c *cliftonDbClient) ConnectContext(ctx context.Context, serverList []string) error {
	if len(serverList) < 1 {
		return EmptyServerListErr
	}
	guessLeader := serverList[0]
	conn, err := grpc.DialContext(ctx, guessLeader, grpc.WithInsecure())

	if err != nil {
		return err
	}
	kvClient := kv_client.NewClientClient(conn)

	rpcCtx, _ := context.WithTimeout(ctx, time.Millisecond*2000)
	res, err := c.kvClient.RegisterClient(
		rpcCtx,
		&kv_client.RegisterClientReq{},
	)

	if err != nil {
		return err
	}

	if res.Status == kv_client.ClientRpcStatus_NOT_LEADER {
		err := conn.Close()
		if err != nil {
			return err
		}
		return c.reconnectWithAddress(res.LeaderHint)
	}
	c.clientId = res.ClientId
	c.leaderAddr = guessLeader
	c.ClientConn = conn
	c.kvClient = kvClient

	return nil
}

func (c *cliftonDbClient) reconnectWithAddress(address *kv_client.Address) error {
	return nil
}

func (c *cliftonDbClient) Connect(serverList []string) error {
	return c.ConnectContext(context.TODO(), serverList)
}

func (c *cliftonDbClient) Close() error {
	err := c.ClientConn.Close()
	c.ClientConn = nil
	c.kvClient = nil

	return err
}

func (c *cliftonDbClient) Get(key string) (value []byte, err error) {
	if c.kvClient != nil {
		return nil, NotConnectedErr
	}
	return nil, nil
}

func (c *cliftonDbClient) Put(key string, value [] []byte) (err error) {
	if c.kvClient != nil {
		return NotConnectedErr
	}
	return nil
}

func (c *cliftonDbClient) Delete(key string) (err error) {
	if c.kvClient != nil {
		return NotConnectedErr
	}
	return nil
}
