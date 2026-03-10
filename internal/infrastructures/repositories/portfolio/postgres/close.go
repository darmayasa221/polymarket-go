package postgres

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

const errClosePositionFailed = "PORTFOLIO.CLOSE_FAILED"

// Close marks a position as closed with the given exit price and close time.
func (r *Repository) Close(ctx context.Context, positionID string, exitPrice decimal.Decimal, closedAt time.Time) error {
	const query = `UPDATE positions SET closed_at = $1, exit_price = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, closedAt, exitPrice.String(), positionID)
	if err != nil {
		return errtypes.NewInternalServerError(errClosePositionFailed)
	}
	return nil
}
