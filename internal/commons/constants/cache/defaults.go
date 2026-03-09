package cache

import "time"

const (
	DefaultUserCacheTTL    = 5 * time.Minute
	DefaultSessionCacheTTL = 24 * time.Hour

	// String forms used as fallback defaults in string-based config parsers.
	DefaultDialTimeout = "5s"
	DefaultReadTimeout = "3s"

	// Duration forms used directly where time.Duration is needed.
	DefaultDialTimeoutDuration = 5 * time.Second
	DefaultReadTimeoutDuration = 3 * time.Second

	// DefaultAddress is the default Redis server address for local/test environments.
	DefaultAddress = "localhost:6379"
)
