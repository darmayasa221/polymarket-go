// Package redis provides the Redis cache adapter.
package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// Client wraps go-redis client.
type Client struct {
	client *goredis.Client
}

// New creates a new Redis client using the given config.
// The startup ping uses context.Background — callers must impose their own deadline via DialTimeout in Config.
func New(cfg Config) (*Client, error) {
	c := goredis.NewClient(&goredis.Options{
		Addr:        cfg.Address,
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}
	return &Client{client: c}, nil
}

// Client returns the underlying go-redis client.
func (c *Client) Client() *goredis.Client { return c.client }

// Ping checks the Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close closes the Redis connection.
func (c *Client) Close() error { return c.client.Close() }
