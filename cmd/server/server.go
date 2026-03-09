package main

import "context"

// Server defines the interface for a runnable server.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}
