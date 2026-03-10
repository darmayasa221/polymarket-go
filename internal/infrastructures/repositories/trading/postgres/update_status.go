package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

const errUpdateOrderStatusFailed = "TRADING.UPDATE_ORDER_STATUS_FAILED"

// UpdateStatus updates the status of the order identified by orderID.
func (r *Repository) UpdateStatus(ctx context.Context, orderID polyid.OrderID, status order.OrderStatus) error {
	const query = `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, string(status), string(orderID))
	if err != nil {
		return errtypes.NewInternalServerError(errUpdateOrderStatusFailed)
	}
	return nil
}
