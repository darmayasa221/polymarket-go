package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

const errListClosedPositionsFailed = "PORTFOLIO.LIST_CLOSED_FAILED"

// ListClosed returns all positions that have been closed.
func (r *Repository) ListClosed(ctx context.Context) ([]*position.Position, error) {
	const query = `SELECT id, asset, token_id, outcome, size, avg_price, market_id, opened_at, closed_at
		FROM positions WHERE closed_at IS NOT NULL`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errListClosedPositionsFailed)
	}
	defer rows.Close()
	return scanPositions(rows)
}
