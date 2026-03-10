package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

const errListOpenOrdersFailed = "TRADING.LIST_OPEN_ORDERS_FAILED"

// ListOpenByMarket returns all open orders for the given market.
func (r *Repository) ListOpenByMarket(ctx context.Context, marketID string) ([]*order.Order, error) {
	const query = `SELECT id, market_id, token_id, side, outcome, price, size, order_type, expiration, fee_rate_bps, signature_type, status, created_at
		FROM orders WHERE market_id = $1 AND status = 'open'`
	rows, err := r.db.QueryContext(ctx, query, marketID)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errListOpenOrdersFailed)
	}
	defer rows.Close()
	return scanOrders(rows)
}
