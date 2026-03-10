package postgres

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

// FindByID returns the position with the given ID.
func (r *Repository) FindByID(ctx context.Context, positionID string) (*position.Position, error) {
	const query = `SELECT id, asset, token_id, outcome, size, avg_price, market_id, opened_at, closed_at
		FROM positions WHERE id = $1`
	return scanPosition(r.db.QueryRowContext(ctx, query, positionID))
}
