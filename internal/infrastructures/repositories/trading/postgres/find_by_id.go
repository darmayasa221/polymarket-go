package postgres

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// FindByID returns the order with the given ID.
func (r *Repository) FindByID(ctx context.Context, orderID polyid.OrderID) (*order.Order, error) {
	const query = `SELECT id, market_id, token_id, side, outcome, price, size, order_type, expiration, fee_rate_bps, signature_type, status, created_at
		FROM orders WHERE id = $1`
	return scanOrder(r.db.QueryRowContext(ctx, query, string(orderID)))
}
