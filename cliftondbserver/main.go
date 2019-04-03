package cliftondbserver

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

func LoadConfig(config *Config) (err error) {

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
	var config Config = Config{
		Server: ApiServer{
			ListenPort: 9091,
		},
		DbPath: "/tmp/cliftondb-test",
		Nodes: RaftNodes{
			SelfId: 1,
			Port:   8080,
			PeerList: []Peer{
				Peer{
					Id:       1,
					IpOrHost: "localhost",
					Port:     8080,
				},
			},
		},
	}

	err := LoadConfig(&config)
	if err != nil {
		os.Exit(-1)
	}

	kvServer, err := NewCliftonDbServer(config)
	log.Printf("creating kvserver")
	if err != nil {
		log.Fatalln("error create new kvserver", err)
	}

	log.Printf("bootstraping kvserver")
	err = kvServer.Boostrap()
	if err != nil {
		log.Fatalln("bootstrap failed", err)
	}

	log.Printf("starting GRPC api server")
	err, doneC := kvServer.ServeGrpc()

	if err != nil {
		log.Fatalln("failed to server grpc", err)
	}
	log.Printf("serving request on localhost:%d", config.Server.ListenPort)
	<-doneC
}
