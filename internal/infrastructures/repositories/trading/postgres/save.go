package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

const errSaveOrderFailed = "TRADING.SAVE_ORDER_FAILED"

// Save persists an order to the database.
// FOK orders have zero Expiration — stored as NULL.
func (r *Repository) Save(ctx context.Context, o *order.Order) error {
	const query = `
		INSERT INTO orders (id, market_id, token_id, side, outcome, price, size, order_type, expiration, fee_rate_bps, signature_type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	var expiration any
	if !o.Expiration().IsZero() {
		expiration = o.Expiration()
	}

	_, err := r.db.ExecContext(ctx, query,
		string(o.ID()),
		o.MarketID(),
		string(o.TokenID()),
		int(o.Side()),
		string(o.Outcome()),
		o.Price().String(),
		o.Size().String(),
		string(o.Type()),
		expiration,
		o.FeeRateBps(),
		o.SignatureType(),
		string(o.Status()),
		o.CreatedAt(),
	)
	if err != nil {
		return errtypes.NewInternalServerError(errSaveOrderFailed)
	}
	return nil
}
