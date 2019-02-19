package config

import "io"

type Config interface {
	Get(key string) string
}

type ConfigReader interface {
	ReadYaml(reader io.Reader) Config
	ReadJson(reader io.Reader) Config
}
