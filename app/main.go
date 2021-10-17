package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	// Construct a server to service the requests against the mux.
	httpServer := http.Server{
		Addr: "localhost:3000",
	}

	// Channel for listening errors from http server
	serverErrors := make(chan error, 1)

	// Starting the server
	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	if err := <-serverErrors; err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
