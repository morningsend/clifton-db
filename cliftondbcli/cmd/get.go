package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zl14917/MastersProject/api/kv-client"
	"google.golang.org/grpc"
	"log"
	"time"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("key", "", "GET request requires a key to retrieve the corresponding value")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Send a GET request to CliftonDB cluster",
	Long:  "Send GET request to CliftonDB cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		serverAddress, err := cmd.Parent().PersistentFlags().GetString("address")
		if err != nil {
			return err
		}

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}
		if len(key) < 1 {
			return fmt.Errorf("key must not be empty")
		}
		executeGet(serverAddress, key)
		return nil
	},
}

func executeGet(serverAddress string, key string) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error connecting to server on %s (%v)", serverAddress, err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("error closing connection %v", err)
		}
	}()

	client := kv_client.NewKVStoreClient(conn)
	timeout := 5 * time.Second
	ctx, _ := context.WithDeadline(
		context.Background(),
		time.Now().Add(timeout),
	)
	req := &kv_client.GetReq{
		Key: key,
	}

	value, err := client.Get(ctx, req)

	fmt.Printf("GET Response Key: %s, Value: %s\n", value.Key, string(value.Value))
}
