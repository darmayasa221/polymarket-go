package redis

import (
	"context"
	"time"
)

// Set stores a value in Redis with a TTL.
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value from Redis. Returns redis.Nil error if key does not exist.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete removes a key from Redis.
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
