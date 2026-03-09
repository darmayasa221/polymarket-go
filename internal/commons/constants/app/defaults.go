package app

import "time"

const (
	// String forms — used as fallback defaults in string-based config parsers.
	DefaultTimeout      = "30s"
	DefaultReadTimeout  = "10s"
	DefaultWriteTimeout = "10s"
	DefaultIdleTimeout  = "60s"

	// Duration forms — used directly where time.Duration is needed.
	DefaultTimeoutDuration      = 30 * time.Second
	DefaultReadTimeoutDuration  = 10 * time.Second
	DefaultWriteTimeoutDuration = 10 * time.Second
	DefaultIdleTimeoutDuration  = 60 * time.Second
	DefaultShutdownTimeout      = 10 * time.Second

	DefaultMaxBodySizeMB = 10
	DefaultLogLevel      = "info"

	// CORS and rate-limit defaults.
	DefaultAllowedOrigins     = "*"
	DefaultRateLimitRequests  = 100
	DefaultRateLimitWindowStr = "1m"
	DefaultRateLimitWindow    = time.Minute
)
