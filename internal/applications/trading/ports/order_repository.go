package ports

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// OrderRepository persists and retrieves Order aggregates.
type OrderRepository interface {
	Save(ctx context.Context, o *order.Order) error
	FindByID(ctx context.Context, orderID polyid.OrderID) (*order.Order, error)
	ListOpenByMarket(ctx context.Context, marketID string) ([]*order.Order, error)
	UpdateStatus(ctx context.Context, orderID polyid.OrderID, status order.OrderStatus) error
}
