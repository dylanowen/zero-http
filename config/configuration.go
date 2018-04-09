package config

import (
	"github.com/dylanowen/zero-http/server"
	"path/filepath"
)

type Configuration struct {
	Server *server.Config

	ConfigDir string
	Debug     bool
}

type RawConfiguration struct {
	Port string

	CertFile string
	KeyFile  string

	Debug bool
}

func NewConfiguration(rawConfiguration RawConfiguration, configDir string) *Configuration {
	return &Configuration{
		Server:    rawConfiguration.NewServerConfig(configDir),
		ConfigDir: configDir,
		Debug:     rawConfiguration.Debug,
	}
}

func (c *RawConfiguration) NewServerConfig(path string) *server.Config {
	var certFile = withPath(path, c.CertFile)
	var keyFile = withPath(path, c.KeyFile)

	return &server.Config{
		Port:     c.Port,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

func withPath(path string, file string) string {
	if len(file) > 0 && !filepath.IsAbs(file) {
		return filepath.Join(path, file)
	}

	return file
}
