// Package main is the entry point for the HTTP server.
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := build()
	if err != nil {
		log.Fatalf("failed to build application: %v", err)
	}

	for _, srv := range app.servers {
		go func(s Server) {
			if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("server stopped unexpectedly: %v", err)
			}
		}(srv)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.shutdown()
}
