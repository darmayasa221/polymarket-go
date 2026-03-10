package ports

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// OrderSubmitter submits signed orders to the Polymarket CLOB.
// The signature is produced by the interfaces layer (private key signing).
type OrderSubmitter interface {
	Submit(ctx context.Context, o *order.Order, signature []byte) (string, error) // returns CLOB orderID
	Cancel(ctx context.Context, clobOrderID string) error
}
