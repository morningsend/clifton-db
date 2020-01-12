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

	log.Printf("starting server with config: %v", config)

	opts := server.DefaultServerOptions(config.RaftID, config.LocalStoragePath)
	dbServer := server.New(opts)

	if err := dbServer.StartUp(); err != nil {
		log.Fatalf("server start up error: %v", err)
	}

	if err := dbServer.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
