package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// GetIDByUsername retrieves only the user's ID by username, checking the Redis cache first.
// Caches the raw ID string, not the full user entity.
// Cache errors are silently ignored — never fail the request for a cache error.
func (r *Repository) GetIDByUsername(ctx context.Context, username string) (user.UserID, error) {
	key := rediscache.UserIDByUsernameKey(username)

	// Cache hit.
	if data, err := r.cache.Get(ctx, key); err == nil && data != "" {
		return user.UserID(data), nil
	}

	// Cache miss — delegate to base.
	id, err := r.base.GetIDByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	// Store in cache (best-effort, ignore error).
	_ = r.cache.Set(ctx, key, id.String(), defaultTTL)

	return id, nil
}
