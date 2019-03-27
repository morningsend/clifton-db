package kvclient

import (
	"context"
	"fmt"
	"github.com/zl14917/MastersProject/api/kv-client"
	"google.golang.org/grpc"
)

type Client struct {
	grpcConn      *grpc.ClientConn
	kvStoreClient kv_client.KVStoreClient
}

func (c *Client) Connect(host string, port uint32) (err error) {
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())

	if err != nil {
		return err
	}

	c.grpcConn = conn
	c.kvStoreClient = kv_client.NewKVStoreClient(conn)

	return nil
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Put(key string, data []byte) (ok bool, err error) {
	result, err := c.kvStoreClient.Put(
		context.Background(),
		&kv_client.PutReq{
			Key:   key,
			Value: data,
		},
	)

	if err != nil {
		return false, err
	}

	return result.Success, nil
}

func (c *Client) Get(key string) (data []byte, err error) {

	value, err := c.kvStoreClient.Get(
		context.Background(),
		&kv_client.GetReq{
			Key: key,
		})

	if err != nil {
		return nil, err
	}

	return value.Value, nil
}
