package main

import (
	"log"

	"github.com/morningsend/clifton-db/pkg/server"
)

var (
	configFile = serverCmd.Flag("config", "Path to config file").Required().ExistingFile()
)

func runServer() {
	config, err := ServerConfigFromTOML(*configFile)
	if err != nil {
		log.Fatalf("error reading config file at %s: %v", *configFile, err)
	}

	opts := server.DefaultServerOptions(config.RaftID, config.LocalStoragePath)
	dbServer := server.New(opts)
	if err := dbServer.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
