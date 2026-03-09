package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
	rediscache "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Add persists a new user via the base repository, then invalidates the username
// existence cache so that VerifyUsername reflects the newly registered user.
// Cache errors are silently ignored — never fail the request for a cache error.
func (r *Repository) Add(ctx context.Context, u *user.User) error {
	if err := r.base.Add(ctx, u); err != nil {
		return err
	}
	_ = r.cache.Delete(ctx, rediscache.UserExistsKey(u.Username()))
	return nil
}
