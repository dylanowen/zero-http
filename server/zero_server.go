package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ZeroServer struct {
	config *Config

	Server *http.Server
}

func NewZeroServer(config *Config) *ZeroServer {
	return &ZeroServer{
		config: config,
	}
}

func (z *ZeroServer) Start() error {
	var handler = http.NewServeMux()

	// get our current working directory
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	handler.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(workingDirectory))))

	var port = z.config.Port

	z.Server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: requestWrapper(handler),
	}

	log.Println("Starting server at localhost" + z.Server.Addr)

	err = z.Server.ListenAndServe()
	// filter out errors from the server closing
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func requestWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		handler.ServeHTTP(w, r)
	})
}

func (z *ZeroServer) Stop() {
	if z.Server != nil {
		var shutdownCtx, _ = context.WithTimeout(context.Background(), 5*time.Second)

		z.Server.Shutdown(shutdownCtx)

		z.Server = nil
	}
}
