// Package order_test tests the order domain.
package order_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

func validOrderParams() order.Params {
	return order.Params{
		MarketID:      "gamma-event-123",
		TokenID:       polyid.TokenID("111"),
		Side:          order.Buy,
		Outcome:       market.Up,
		Price:         decimal.NewFromFloat(0.65),
		Size:          decimal.NewFromFloat(10),
		Type:          order.GTD,
		Expiration:    timeutil.Now().Add(5 * time.Minute),
		FeeRateBps:    50,
		SignatureType: 0,
	}
}

func TestNew_Valid(t *testing.T) {
	t.Parallel()

	o, err := order.New(validOrderParams())
	require.NoError(t, err)
	assert.Equal(t, order.Buy, o.Side())
	assert.Equal(t, order.StatusOpen, o.Status())
	assert.False(t, o.ID().IsEmpty())
}

func TestNew_MissingMarketID(t *testing.T) {
	t.Parallel()

	p := validOrderParams()
	p.MarketID = ""
	_, err := order.New(p)
	assert.ErrorContains(t, err, order.ErrMarketIDRequired)
}

func TestNew_InvalidPrice(t *testing.T) {
	t.Parallel()

	p := validOrderParams()
	p.Price = decimal.Zero
	_, err := order.New(p)
	assert.ErrorContains(t, err, order.ErrPriceInvalid)
}

func TestNew_InvalidSize(t *testing.T) {
	t.Parallel()

	p := validOrderParams()
	p.Size = decimal.Zero
	_, err := order.New(p)
	assert.ErrorContains(t, err, order.ErrSizeInvalid)
}

func TestNew_GTDRequiresExpiration(t *testing.T) {
	t.Parallel()

	p := validOrderParams()
	p.Type = order.GTD
	p.Expiration = time.Time{}
	_, err := order.New(p)
	assert.ErrorContains(t, err, order.ErrExpirationRequired)
}
