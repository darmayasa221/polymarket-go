package cached

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// ListOffset delegates directly to the base repository without caching.
// List queries are not cached: invalidation on any write would be too complex.
func (r *Repository) ListOffset(ctx context.Context, params pagination.OffsetParams) (pagination.OffsetResult[*user.User], error) {
	return r.base.ListOffset(ctx, params)
}
