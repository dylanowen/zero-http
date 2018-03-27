package config

import (
	"github.com/dylanowen/zero-http/server"
)

type Configuration struct {
	RawConfiguration

	ConfigDir string
}

type RawConfiguration struct {
	Server server.Config
	Debug  bool
}

func getDefault() *RawConfiguration {
	return &RawConfiguration{
		Server: *server.DefaultConfig(),
		Debug:  false,
	}
}

func appendPath(path string, key string) string {
	var currentPath = path
	if len(path) > 0 {
		currentPath += "."
	}
	currentPath += key

	return currentPath
}
