package server

import (
	"context"
	"github.com/phayes/freeport"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ZeroServer struct {
	config *Config

	Server *http.Server
}

func NewServer(config *Config) *ZeroServer {
	return &ZeroServer{
		config: config,
	}
}

func (z *ZeroServer) Start() (err error) {
	// find our port
	var port int
	if port, err = z.findPort(); err != nil {
		// if we didn't find a port return an error
		return err
	}

	// setup our handler
	var handler *http.ServeMux
	if handler, err = getHandler(); err != nil {
		return err
	}

	// create our server
	z.Server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: requestWrapper(handler),
	}

	// start the actual server
	err = z.listenAndServe()

	// filter out errors from the server closing
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (z *ZeroServer) Stop() {
	if z.Server != nil {
		var shutdownCtx, _ = context.WithTimeout(context.Background(), 5*time.Second)

		z.Server.Shutdown(shutdownCtx)

		z.Server = nil
	}
}

func requestWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// Don't cache anything
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		handler.ServeHTTP(w, r)
	})
}

func (z *ZeroServer) findPort() (port int, err error) {
	var rawPort = strings.TrimSpace(z.config.Port)
	if "any" == strings.ToLower(rawPort) {
		port, err = freeport.GetFreePort()
	} else if len(rawPort) <= 0 || rawPort == "0" {
		log.Println("Finding any port since none were entered")

		port, err = freeport.GetFreePort()
	} else {
		var parsedPort int64
		parsedPort, err = strconv.ParseInt(rawPort, 0, 32)
		port = int(parsedPort)
	}

	return
}

func getHandler() (*http.ServeMux, error) {
	// get our current working directory
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var handler = http.NewServeMux()

	handler.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(workingDirectory))))

	return handler, nil
}

func (z *ZeroServer) listenAndServe() (err error) {
	// check to see if this is a TLS server
	var certFile = z.config.CertFile
	var keyFile = z.config.KeyFile
	if len(certFile) > 0 && len(keyFile) > 0 {
		// check for missing files
		_, certErr := os.Stat(certFile)
		if certErr != nil {
			log.Println("Couldn't load certFile:", certFile)
		}
		_, keyErr := os.Stat(keyFile)
		if keyErr != nil {
			log.Println("Couldn't load keyFile:", keyFile)
		}

		// if there are no errors we have our certificates
		if certErr == nil && keyErr == nil {
			err = z.serveHttps(certFile, keyFile)
		} else {
			err = z.serverHttp()
		}
	} else {
		err = z.serverHttp()
	}

	return
}

func (z *ZeroServer) serverHttp() error {
	log.Println("Starting server at https://localhost" + z.Server.Addr)

	return z.Server.ListenAndServe()
}

func (z *ZeroServer) serveHttps(certFile string, keyFile string) error {
	log.Println("Starting server at http://localhost" + z.Server.Addr)

	return z.Server.ListenAndServeTLS(certFile, keyFile)
}
