// Package position_test tests the position domain.
package position_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

func validPositionParams() position.Params {
	return position.Params{
		Asset:    market.BTC,
		TokenID:  polyid.TokenID("111"),
		Outcome:  market.Up,
		Size:     decimal.NewFromFloat(10),
		AvgPrice: decimal.NewFromFloat(0.65),
		MarketID: "gamma-event-123",
	}
}

func TestNew_Valid(t *testing.T) {
	t.Parallel()

	pos, err := position.New(validPositionParams())
	require.NoError(t, err)
	assert.Equal(t, market.BTC, pos.Asset())
	assert.False(t, pos.ID() == "")
	assert.Nil(t, pos.ClosedAt())
}

func TestNew_MissingMarketID(t *testing.T) {
	t.Parallel()

	p := validPositionParams()
	p.MarketID = ""
	_, err := position.New(p)
	assert.ErrorContains(t, err, position.ErrMarketIDRequired)
}

func TestNew_InvalidSize(t *testing.T) {
	t.Parallel()

	p := validPositionParams()
	p.Size = decimal.Zero
	_, err := position.New(p)
	assert.ErrorContains(t, err, position.ErrSizeInvalid)
}

func TestNew_InvalidAvgPrice(t *testing.T) {
	t.Parallel()

	p := validPositionParams()
	p.AvgPrice = decimal.Zero
	_, err := position.New(p)
	assert.ErrorContains(t, err, position.ErrAvgPriceInvalid)
}
