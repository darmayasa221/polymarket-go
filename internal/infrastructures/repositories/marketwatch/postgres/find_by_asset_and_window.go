package postgres

import (
	"context"
	"time"

	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// FindByAssetAndWindow returns the market for the given asset and window start time.
func (r *Repository) FindByAssetAndWindow(ctx context.Context, asset market.Asset, windowStart time.Time) (*market.Market, error) {
	const query = `SELECT id, slug, asset, window_start, condition_id, up_token_id, down_token_id, tick_size, fee_enabled, active
		FROM markets WHERE asset = $1 AND window_start = $2 LIMIT 1`
	return scanMarket(r.db.QueryRowContext(ctx, query, string(asset), windowStart))
}
