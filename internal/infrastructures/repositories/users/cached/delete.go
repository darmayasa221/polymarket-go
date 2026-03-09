package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Delete removes a user by ID via the base repository,
// then invalidates all related cache keys.
// Cache errors are silently ignored — never fail the request for a cache error.
func (r *Repository) Delete(ctx context.Context, id user.UserID) error {
	// Fetch from base (not cache) to retrieve username for key invalidation.
	// This is best-effort — if it fails, we still proceed with the delete.
	u, _ := r.base.GetByID(ctx, id)

	if err := r.base.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate ID-keyed cache entry.
	_ = r.cache.InvalidateUser(ctx, id.String())

	// Invalidate username-keyed cache entries if we resolved the username.
	if u != nil {
		_ = r.cache.Delete(ctx, rediscache.UserByUsernameKey(u.Username()))
		_ = r.cache.Delete(ctx, rediscache.UserIDByUsernameKey(u.Username()))
		_ = r.cache.Delete(ctx, rediscache.UserExistsKey(u.Username()))
	}

	return nil
}
