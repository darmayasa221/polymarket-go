package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

const errListOpenPositionsFailed = "PORTFOLIO.LIST_OPEN_FAILED"

// ListOpen returns all positions that have not been closed.
func (r *Repository) ListOpen(ctx context.Context) ([]*position.Position, error) {
	const query = `SELECT id, asset, token_id, outcome, size, avg_price, market_id, opened_at, closed_at
		FROM positions WHERE closed_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errListOpenPositionsFailed)
	}
	defer rows.Close()
	return scanPositions(rows)
}
