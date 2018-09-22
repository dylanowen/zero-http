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
	if handler, err = z.getHandler(); err != nil {
		return err
	}

	// create our server
	z.Server = &http.Server{
		Addr:    z.getHost() + ":" + strconv.Itoa(port),
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

		if !handleOptions(w, r) {
			// Don't cache anything
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			var origin = r.Header.Get("Origin")
			if origin != "" {
				// auto allow everything for CORS
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			handler.ServeHTTP(w, r)
		}
	})
}

const allowedMethods = "GET,OPTIONS"

func handleOptions(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "OPTIONS" {
		var origin = r.Header.Get("Origin")
		var requestMethod = r.Header.Get("Access-Control-Request-Method")

		// This is a CORS request so respond appropriately
		if origin != "" && requestMethod != "" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Origin", origin)

			var requestHeaders = r.Header.Get("Access-Control-Request-Headers")
			if requestHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
			}
		}

		w.Header().Set("Allow", allowedMethods)

		return true
	} else {
		return false
	}
}

func (z *ZeroServer) getHost() string {
	var host = z.config.Host

	// make the user explicitly choose 0.0.0.0 instead of auto binding
	if host == "" {
		host = "localhost"
	}

	return host
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

func (z *ZeroServer) getHandler() (*http.ServeMux, error) {
	var prefix = z.config.BasePath
	if prefix == "" {
		prefix = "/"
	}

	// print some base path warnings
	if prefix[0] != '/' {
		log.Println("Expected leading slash for base path")
	}
	if len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
		log.Println("A trailing base path slash could cause unexpected results")
	}

	// get our current working directory
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var handler = http.NewServeMux()

	handler.Handle("/", http.StripPrefix(prefix, http.FileServer(http.Dir(workingDirectory))))

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
	log.Println("Starting server at http://" + z.Server.Addr)

	return z.Server.ListenAndServe()
}

func (z *ZeroServer) serveHttps(certFile string, keyFile string) error {
	log.Println("Starting server at https://" + z.Server.Addr)

	return z.Server.ListenAndServeTLS(certFile, keyFile)
}
