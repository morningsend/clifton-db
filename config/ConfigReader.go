package config

import "io"

type Config interface {
	Get(key string) string
}

type Reader interface {
	ReadYaml(reader io.Reader) Config
	ReadJson(reader io.Reader) Config
}
