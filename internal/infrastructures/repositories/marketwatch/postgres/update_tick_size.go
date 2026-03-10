package postgres

import (
	"context"

	"github.com/shopspring/decimal"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

const errUpdateTickSizeFailed = "MARKETWATCH.UPDATE_TICK_SIZE_FAILED"

// UpdateTickSize updates the tick size for a market identified by its condition ID.
func (r *Repository) UpdateTickSize(ctx context.Context, conditionID string, newTickSize decimal.Decimal) error {
	const query = `UPDATE markets SET tick_size = $1 WHERE condition_id = $2`
	_, err := r.db.ExecContext(ctx, query, newTickSize.String(), conditionID)
	if err != nil {
		return errtypes.NewInternalServerError(errUpdateTickSizeFailed)
	}
	return nil
}
