package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.PersistentFlags().StringArray("address", []string{}, "GRPC server address including host and port [host]:[port]")
}

var rootCmd = &cobra.Command{
	Use:   "cliftondbcli",
	Short: "CliftonDB CLI provides functionality on command line to db cluster",
	Long:  "CliftonDB CLI tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
