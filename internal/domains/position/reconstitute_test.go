package position_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

func TestReconstitute_Position(t *testing.T) {
	t.Parallel()
	openedAt := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	closedAt := time.Date(2026, 3, 10, 12, 4, 0, 0, time.UTC)
	pos := position.Reconstitute(position.ReconstitutedParams{
		ID:       "pos-uuid",
		Asset:    market.BTC,
		TokenID:  polyid.TokenID("token-up"),
		Outcome:  market.Up,
		Size:     decimal.RequireFromString("10"),
		AvgPrice: decimal.RequireFromString("0.60"),
		MarketID: "market-1",
		OpenedAt: openedAt,
		ClosedAt: &closedAt,
	})
	require.NotNil(t, pos)
	assert.Equal(t, "pos-uuid", pos.ID())
	assert.NotNil(t, pos.ClosedAt())
}
