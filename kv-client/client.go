package kv_client

import (
	"context"
	"errors"
	"github.com/zl14917/MastersProject/api/kv-client"
	"google.golang.org/grpc"
	"time"
)

var NotConnectedErr = errors.New("client not connected to server")
var PutFailedErr = errors.New("put key value failed")
var DeleteFailedErr = errors.New("delete failed")

type Client struct {
	ServerAddress string
	conn          *grpc.ClientConn
	kvStoreClient kv_client.KVStoreClient
	conf          Config
}

type Config struct {
	Timeout       time.Duration
	ServerAddress string
}

func NewClient(conf Config) *Client {
	return &Client{
		kvStoreClient: nil,
		conn:          nil,
		conf:          conf,
	}
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(c.conf.ServerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.kvStoreClient = kv_client.NewKVStoreClient(conn)
	return nil
}
func (c *Client) Get(key string) ([]byte, error) {
	if c.kvStoreClient == nil {
		return nil, NotConnectedErr
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		c.conf.Timeout,
	)
	req := kv_client.GetReq{
		Key: key,
	}
	val, err := c.kvStoreClient.Get(ctx, &req)
	if err != nil {
		return nil, err
	}
	return val.Value, nil
}

func (c *Client) Put(key string, value []byte) error {
	if c.kvStoreClient == nil {
		return NotConnectedErr
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		c.conf.Timeout,
	)
	req := kv_client.PutReq{
		Key:   key,
		Value: value,
	}
	res, err := c.kvStoreClient.Put(ctx, &req)

	if err != nil {
		return err
	}

	if res.Success {
		return nil
	}
	return PutFailedErr
}

func (c *Client) Delete(key string) error {
	if c.kvStoreClient == nil {
		return NotConnectedErr
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		c.conf.Timeout,
	)

	req := kv_client.DelReq{
		Key: key,
	}

	res, err := c.kvStoreClient.Delete(ctx, &req)
	if err != nil {
		return err
	}
	if res.Ok {
		return nil
	}
	return DeleteFailedErr
}
