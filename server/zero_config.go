package server

type Config struct {
	Port int
}

func DefaultConfig() *Config {
	return &Config{
		Port: 8000,
	}
}
