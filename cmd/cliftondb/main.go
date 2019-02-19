package cliftondb

import (
	"flag"
	"fmt"
	"github.com/zl14917/MastersProject/kvserver"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
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

func main() {
	flag.Parse()
	if *configPath == "" {
		fmt.Fprintf(os.Stderr, "Error reading config path, missing.")
		os.Exit(-1)
	}
	absPath, err := resolveAbsConfigPath(*configPath)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error resolving config file path: %s: %v\n", *configPath, err)
		os.Exit(-1)
	}

	file, err := os.Open(absPath)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error opening config file: absolute path: %s: %v\n", absPath, err)
		os.Exit(-1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error close config file: %v\n", err)
		}
	}()

	var config kvserver.Config
	yamlDecoder := yaml.NewDecoder(file)
	err = yamlDecoder.Decode(&config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
	}
}
