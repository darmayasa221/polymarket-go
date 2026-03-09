package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// GetByID retrieves a user by ID, checking the Redis cache first.
// On cache miss, delegates to the base repository and stores the result.
// Cache errors are silently ignored — never fail the request for a cache error.
func (r *Repository) GetByID(ctx context.Context, id user.UserID) (*user.User, error) {
	key := rediscache.UserKey(id.String())

	// Cache hit.
	if data, err := r.cache.Get(ctx, key); err == nil {
		if u := unmarshalUser(data); u != nil {
			return u, nil
		}
	}

	// Cache miss — delegate to base.
	u, err := r.base.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache (best-effort, ignore error).
	if data, ok := marshalUser(u); ok {
		_ = r.cache.Set(ctx, key, data, defaultTTL)
	}

	return u, nil
}
