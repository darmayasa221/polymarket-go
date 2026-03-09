// Package oracle_test tests the oracle domain.
package oracle_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

func validPriceParams() oracle.Params {
	return oracle.Params{
		Asset:      "btc",
		Source:     oracle.SourceChainlink,
		Value:      decimal.NewFromFloat(65000.50),
		RoundedAt:  timeutil.Now().Add(-2),
		ReceivedAt: timeutil.Now(),
	}
}

func TestNewPrice_Valid(t *testing.T) {
	t.Parallel()

	p, err := oracle.New(validPriceParams())
	require.NoError(t, err)
	assert.Equal(t, "btc", p.Asset())
	assert.Equal(t, oracle.SourceChainlink, p.Source())
}

func TestNewPrice_MissingAsset(t *testing.T) {
	t.Parallel()

	params := validPriceParams()
	params.Asset = ""
	_, err := oracle.New(params)
	assert.ErrorContains(t, err, oracle.ErrAssetRequired)
}

func TestNewPrice_ZeroValue(t *testing.T) {
	t.Parallel()

	params := validPriceParams()
	params.Value = decimal.Zero
	_, err := oracle.New(params)
	assert.ErrorContains(t, err, oracle.ErrPriceValueInvalid)
}

func TestNewPrice_InvalidSource(t *testing.T) {
	t.Parallel()

	params := validPriceParams()
	params.Source = oracle.PriceSource("unknown")
	_, err := oracle.New(params)
	assert.ErrorContains(t, err, oracle.ErrInvalidSource)
}
