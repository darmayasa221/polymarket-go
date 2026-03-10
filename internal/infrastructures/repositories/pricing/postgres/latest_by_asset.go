package postgres

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// LatestByAsset returns the most recently received price for asset (any source).
func (r *Repository) LatestByAsset(ctx context.Context, asset string) (*oracle.Price, error) {
	const query = `SELECT asset, source, value, rounded_at, received_at
		FROM prices WHERE asset = $1 ORDER BY received_at DESC LIMIT 1`
	return scanPrice(r.db.QueryRowContext(ctx, query, asset))
}
