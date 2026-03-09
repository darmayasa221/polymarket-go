package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Update persists changes to an existing user via the base repository,
// then invalidates all related cache keys.
// Cache errors are silently ignored — never fail the request for a cache error.
func (r *Repository) Update(ctx context.Context, u *user.User) error {
	// Fetch current state to capture the old username before updating.
	old, _ := r.base.GetByID(ctx, u.ID()) // best-effort — ignore error

	if err := r.base.Update(ctx, u); err != nil {
		return err
	}

	// Invalidate by ID.
	_ = r.cache.InvalidateUser(ctx, u.ID().String())

	// Invalidate new username keys.
	_ = r.cache.Delete(ctx, rediscache.UserByUsernameKey(u.Username()))
	_ = r.cache.Delete(ctx, rediscache.UserIDByUsernameKey(u.Username()))
	_ = r.cache.Delete(ctx, rediscache.UserExistsKey(u.Username()))

	// Invalidate old username keys if username changed.
	if old != nil && old.Username() != u.Username() {
		_ = r.cache.Delete(ctx, rediscache.UserByUsernameKey(old.Username()))
		_ = r.cache.Delete(ctx, rediscache.UserIDByUsernameKey(old.Username()))
		_ = r.cache.Delete(ctx, rediscache.UserExistsKey(old.Username()))
	}

	return nil
}
