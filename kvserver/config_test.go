package kvserver

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
)

const fileContent = `
data:
  path-dir: /var/cliftondb/data
log:
  path-dir: /var/cliftondb/logs
  segment-size: 32MB
nodes:
  self-id: 1
  port: 10030
  peers:
    - id: 2
      host: localhost
      port: 10031
    - id: 3
      host: localhost
      port: 10032
`

func TestDeserializeConfigFile(t *testing.T) {
	reader := strings.NewReader(fileContent)
	decoder := yaml.NewDecoder(reader)
	config := Config{}
	err := decoder.Decode(&config)

	if err != nil {
		t.Errorf("error reading Conf file: %v\n", err)
	}
	fmt.Println(config)
}
