package position_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

func TestUnrealisedPnL(t *testing.T) {
	t.Parallel()

	// Bought 10 shares at 0.65, current price 0.75 → unrealised = 10 * (0.75 - 0.65) = 1.00
	pos, err := position.New(position.Params{
		Asset:    "btc",
		TokenID:  "111",
		Outcome:  "Up",
		Size:     decimal.NewFromFloat(10),
		AvgPrice: decimal.NewFromFloat(0.65),
		MarketID: "market-1",
	})
	require.NoError(t, err)

	pnl := pos.UnrealisedPnL(decimal.NewFromFloat(0.75))
	assert.True(t, pnl.Equal(decimal.NewFromFloat(1.00)), "got: %s", pnl.String())
}

func TestRealisedPnL(t *testing.T) {
	t.Parallel()

	// Bought 10 shares at 0.65, sold at 0.80 → realized = 10 * (0.80 - 0.65) = 1.50
	pos, err := position.New(position.Params{
		Asset:    "btc",
		TokenID:  "111",
		Outcome:  "Up",
		Size:     decimal.NewFromFloat(10),
		AvgPrice: decimal.NewFromFloat(0.65),
		MarketID: "market-1",
	})
	require.NoError(t, err)

	pnl := pos.RealisedPnL(decimal.NewFromFloat(0.80))
	assert.True(t, pnl.Equal(decimal.NewFromFloat(1.50)), "got: %s", pnl.String())
}
