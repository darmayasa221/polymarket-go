package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// GetPassword retrieves the hashed password for a username.
// Passwords are security-sensitive — this method never caches.
func (r *Repository) GetPassword(ctx context.Context, username string) (user.HashedPassword, error) {
	return r.base.GetPassword(ctx, username)
}
