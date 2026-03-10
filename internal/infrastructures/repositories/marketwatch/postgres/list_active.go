package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

const errListActiveMarketsFailed = "MARKETWATCH.LIST_ACTIVE_FAILED"

// ListActive returns all markets currently marked as active.
func (r *Repository) ListActive(ctx context.Context) ([]*market.Market, error) {
	const query = `SELECT id, slug, asset, window_start, condition_id, up_token_id, down_token_id, tick_size, fee_enabled, active
		FROM markets WHERE active = true`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errListActiveMarketsFailed)
	}
	defer rows.Close()
	return scanMarkets(rows)
}
