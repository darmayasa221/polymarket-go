// Package server provides the HTTP server configuration for the infrastructures layer.
package server

import "time"

// Config holds the configuration for the HTTP server.
type Config struct {
	// Port is the TCP port the HTTP server listens on (e.g. "3000").
	Port string
	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration
	// IdleTimeout is the maximum amount of time to wait for the next request on a keep-alive connection.
	IdleTimeout time.Duration
	// RequestTimeout is the maximum duration a handler may run before the context is canceled.
	RequestTimeout time.Duration
	// MaxBodyBytes is the maximum number of bytes the server will read from the request body.
	MaxBodyBytes int64
	// AllowedOrigins is the list of origins permitted by the CORS middleware (e.g. ["*"] or ["https://example.com"]).
	AllowedOrigins []string
	// RateLimitRequests is the maximum number of requests allowed per IP per RateLimitWindow.
	RateLimitRequests int
	// RateLimitWindow is the sliding window duration for rate limiting.
	RateLimitWindow time.Duration
}
