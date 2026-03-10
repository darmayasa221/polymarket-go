package market_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

func TestReconstitute_Market(t *testing.T) {
	t.Parallel()
	windowStart := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	m := market.Reconstitute(market.ReconstitutedParams{
		ID:          "gamma-event-id",
		Asset:       market.BTC,
		WindowStart: windowStart,
		ConditionID: polyid.ConditionID("cond-1"),
		UpTokenID:   polyid.TokenID("tok-up"),
		DownTokenID: polyid.TokenID("tok-down"),
		TickSize:    decimal.NewFromFloat(0.01),
		FeeEnabled:  true,
		Active:      false, // can reconstitute inactive market
	})
	require.NotNil(t, m)
	assert.Equal(t, "gamma-event-id", m.ID())
	assert.Equal(t, market.BTC, m.Asset())
	assert.False(t, m.Active())
}
