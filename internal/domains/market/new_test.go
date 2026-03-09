// Package market_test tests the market domain.
package market_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

func validParams() market.Params {
	now := timeutil.WindowStart(time.Unix(1_700_000_100, 0).UTC())
	return market.Params{
		ID:          "gamma-event-123",
		Asset:       market.BTC,
		WindowStart: now,
		ConditionID: polyid.ConditionID("0xabc123"),
		UpTokenID:   polyid.TokenID("111"),
		DownTokenID: polyid.TokenID("222"),
		TickSize:    decimal.NewFromFloat(0.01),
		FeeEnabled:  true,
	}
}

func TestNew_Valid(t *testing.T) {
	t.Parallel()

	m, err := market.New(validParams())
	require.NoError(t, err)
	assert.Equal(t, market.BTC, m.Asset())
	assert.Equal(t, polyid.ConditionID("0xabc123"), m.ConditionID())
	assert.True(t, m.Active())
}

func TestNew_MissingID(t *testing.T) {
	t.Parallel()

	p := validParams()
	p.ID = ""
	_, err := market.New(p)
	assert.ErrorContains(t, err, market.ErrIDRequired)
}

func TestNew_InvalidAsset(t *testing.T) {
	t.Parallel()

	p := validParams()
	p.Asset = market.Asset("invalid")
	_, err := market.New(p)
	assert.ErrorContains(t, err, market.ErrInvalidAsset)
}

func TestNew_MissingConditionID(t *testing.T) {
	t.Parallel()

	p := validParams()
	p.ConditionID = polyid.ConditionID("")
	_, err := market.New(p)
	assert.ErrorContains(t, err, market.ErrConditionIDRequired)
}

func TestNew_MissingTokenIDs(t *testing.T) {
	t.Parallel()

	p := validParams()
	p.UpTokenID = polyid.TokenID("")
	_, err := market.New(p)
	assert.ErrorContains(t, err, market.ErrUpTokenRequired)

	p2 := validParams()
	p2.DownTokenID = polyid.TokenID("")
	_, err2 := market.New(p2)
	assert.ErrorContains(t, err2, market.ErrDownTokenRequired)
}

func TestNew_ZeroWindowStart(t *testing.T) {
	t.Parallel()

	p := validParams()
	p.WindowStart = time.Time{}
	_, err := market.New(p)
	assert.ErrorContains(t, err, market.ErrWindowStartRequired)
}

func TestMarket_WindowEnd(t *testing.T) {
	t.Parallel()

	m, err := market.New(validParams())
	require.NoError(t, err)
	expected := m.WindowStart().Add(5 * time.Minute)
	assert.Equal(t, expected, m.WindowEnd())
}

func TestMarket_Slug(t *testing.T) {
	t.Parallel()

	m, err := market.New(validParams())
	require.NoError(t, err)
	assert.Equal(t, "btc-updown-5m-1700000100", m.Slug().String())
}
