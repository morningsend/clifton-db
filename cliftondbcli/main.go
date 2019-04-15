package main

import "github.com/zl14917/MastersProject/cliftondbcli/cmd"

// Command line tools for interacting with a Clifton Db cluster

var serverAddress = []string{"localhost:9091"}

func main() {
	cmd.Execute()
}
