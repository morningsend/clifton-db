package config

import "io"

type Config interface {
	Get(key string) string
}

type ConfigReader interface {
	Read(reader io.Reader) Config
}

