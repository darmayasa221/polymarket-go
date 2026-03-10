package postgres

import (
	"context"
	"time"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// WindowOpenPrice returns the first Chainlink price at or after windowStart for asset.
func (r *Repository) WindowOpenPrice(ctx context.Context, asset string, windowStart time.Time) (*oracle.Price, error) {
	const query = `SELECT asset, source, value, rounded_at, received_at
		FROM prices WHERE asset = $1 AND source = $2 AND received_at >= $3
		ORDER BY received_at ASC LIMIT 1`
	return scanPrice(r.db.QueryRowContext(ctx, query, asset, string(oracle.SourceChainlink), windowStart))
}
