package cached

import (
	"context"
	"fmt"

	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// VerifyUsername returns true if the username already exists, checking the Redis cache first.
// Caches the boolean result as "true" or "false".
// Cache errors are silently ignored — never fail the request for a cache error.
func (r *Repository) VerifyUsername(ctx context.Context, username string) (bool, error) {
	key := rediscache.UserExistsKey(username)

	// Cache hit.
	if data, err := r.cache.Get(ctx, key); err == nil {
		return data == "true", nil
	}

	// Cache miss — delegate to base.
	exists, err := r.base.VerifyUsername(ctx, username)
	if err != nil {
		return false, err
	}

	// Store in cache (best-effort, ignore error).
	_ = r.cache.Set(ctx, key, fmt.Sprintf("%v", exists), defaultTTL)

	return exists, nil
}
