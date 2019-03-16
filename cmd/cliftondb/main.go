package main

import (
	"flag"
	"fmt"
	"github.com/zl14917/MastersProject/kvserver"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {

}

var configPath = flag.String("config", "", "Path to clifton DB configuration file.")

func resolveAbsConfigPath(path string) (string, error) {

	if path[0] == os.PathSeparator {
		return path, nil
	}

	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, path), nil
}

const fileContent = `
data:
  path-dir: /var/cliftondb/data
log:
  path-dir: /var/cliftondb/logs
  segment-size: 32MB
nodes:
  self-id: 1
  port: 10030
  peers: []
`

const TESTING = true

func LoadConfig(config *kvserver.Config) (err error) {

	if TESTING {
		reader := strings.NewReader(fileContent)
		yamlDecoder := yaml.NewDecoder(reader)
		err := yamlDecoder.Decode(config)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			return err
		}
		return nil
	}

	flag.Parse()
	if *configPath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading config path, missing.")
		return err
	}
	absPath, err := resolveAbsConfigPath(*configPath)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error resolving config file path: %s: %v\n", *configPath, err)
		return err
	}

	file, err := os.Open(absPath)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error opening config file: absolute path: %s: %v\n", absPath, err)
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error close config file: %v\n", err)
		}
	}()

	yamlDecoder := yaml.NewDecoder(file)
	err = yamlDecoder.Decode(config)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		return err
	}
	return nil
}

func main() {
	var config kvserver.Config
	err := LoadConfig(&config)
	if err != nil {
		os.Exit(-1)
	}

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(config.Nodes.Port)))

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot listen on port", config.Nodes.Port)
		os.Exit(-1)
	}

	kvServer, err := kvserver.NewKVServer(config)
	if err != nil {
		log.Fatalln("error create new kvserver", err)
	}

	log.Println("Starting grpc server on port:", config.Nodes.Port)
	grpcServer := grpc.NewServer()
	kvServer.RegisterGrpcServer(grpcServer)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("error serving grpc: %v\n", err)
	}
}
