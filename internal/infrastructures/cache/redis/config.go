// Package redis provides the Redis cache adapter for the infrastructures layer.
package redis

import "time"

// Config holds the configuration for the Redis cache adapter.
type Config struct {
	// Address is the Redis server address (host:port).
	Address string
	// Password is the optional Redis authentication password.
	Password string
	// DB is the Redis database index to use.
	DB int
	// DialTimeout is the timeout for establishing a new connection.
	DialTimeout time.Duration
	// ReadTimeout is the timeout for reading a reply from Redis.
	ReadTimeout time.Duration
}
