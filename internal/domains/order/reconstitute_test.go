package order_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

func TestReconstitute_Order(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	o := order.Reconstitute(order.ReconstitutedParams{
		ID:            polyid.OrderID("order-uuid"),
		MarketID:      "market-1",
		TokenID:       polyid.TokenID("token-1"),
		Side:          order.Buy,
		Outcome:       market.Up,
		Price:         decimal.RequireFromString("0.60"),
		Size:          decimal.RequireFromString("10"),
		Type:          order.FOK,
		FeeRateBps:    100,
		SignatureType: 0,
		Status:        order.StatusFilled,
		CreatedAt:     createdAt,
	})
	require.NotNil(t, o)
	assert.Equal(t, polyid.OrderID("order-uuid"), o.ID())
	assert.Equal(t, order.StatusFilled, o.Status())
}
