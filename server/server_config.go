package server

type Config struct {
	Port int

	CertFile string
	KeyFile  string
}

func DefaultConfig() *Config {
	return &Config{
		Port: 8000,
	}
}
