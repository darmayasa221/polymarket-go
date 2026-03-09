package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// ListCursor delegates directly to the base repository without caching.
// List queries are not cached: invalidation on any write would be too complex.
func (r *Repository) ListCursor(ctx context.Context, params pagination.CursorParams) (pagination.CursorResult[*user.User], error) {
	return r.base.ListCursor(ctx, params)
}
