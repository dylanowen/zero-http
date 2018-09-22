package server

type Config struct {
	Host     string
	Port     string
	BasePath string

	CertFile string
	KeyFile  string
}
