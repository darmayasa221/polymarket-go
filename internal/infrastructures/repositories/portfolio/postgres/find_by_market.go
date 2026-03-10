package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

const errFindByMarketFailed = "PORTFOLIO.FIND_BY_MARKET_FAILED"

// FindByMarket returns all positions for the given market.
func (r *Repository) FindByMarket(ctx context.Context, marketID string) ([]*position.Position, error) {
	const query = `SELECT id, asset, token_id, outcome, size, avg_price, market_id, opened_at, closed_at
		FROM positions WHERE market_id = $1`
	rows, err := r.db.QueryContext(ctx, query, marketID)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errFindByMarketFailed)
	}
	defer rows.Close()
	return scanPositions(rows)
}
