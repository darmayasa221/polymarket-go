package postgres

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// LatestChainlinkByAsset returns the most recently received Chainlink price for asset.
func (r *Repository) LatestChainlinkByAsset(ctx context.Context, asset string) (*oracle.Price, error) {
	const query = `SELECT asset, source, value, rounded_at, received_at
		FROM prices WHERE asset = $1 AND source = $2 ORDER BY received_at DESC LIMIT 1`
	return scanPrice(r.db.QueryRowContext(ctx, query, asset, string(oracle.SourceChainlink)))
}
