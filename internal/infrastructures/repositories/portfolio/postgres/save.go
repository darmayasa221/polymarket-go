package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

const errSavePositionFailed = "PORTFOLIO.SAVE_FAILED"

// Save persists a new position. closed_at and exit_price are NULL for open positions.
func (r *Repository) Save(ctx context.Context, pos *position.Position) error {
	const query = `
		INSERT INTO positions (id, asset, token_id, outcome, size, avg_price, market_id, opened_at, closed_at, exit_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL, NULL)`

	_, err := r.db.ExecContext(ctx, query,
		pos.ID(),
		string(pos.Asset()),
		string(pos.TokenID()),
		string(pos.Outcome()),
		pos.Size().String(),
		pos.AvgPrice().String(),
		pos.MarketID(),
		pos.OpenedAt(),
	)
	if err != nil {
		return errtypes.NewInternalServerError(errSavePositionFailed)
	}
	return nil
}
