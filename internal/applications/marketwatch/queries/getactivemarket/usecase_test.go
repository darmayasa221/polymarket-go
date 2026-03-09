package getactivemarket_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// errMockNotFound satisfies nilnil linter for unused mock methods.
var errMockNotFound = errors.New("mock: not found")

type mockMarketRepo struct {
	found    *market.Market
	foundErr error
}

func (m *mockMarketRepo) Save(_ context.Context, _ *market.Market) error {
	return errors.New("mock: Save not used")
}
func (m *mockMarketRepo) FindByAssetAndWindow(_ context.Context, _ market.Asset, _ time.Time) (*market.Market, error) {
	if m.foundErr != nil {
		return nil, m.foundErr
	}
	if m.found != nil {
		return m.found, nil
	}
	return nil, errMockNotFound
}
func (m *mockMarketRepo) UpdateTickSize(_ context.Context, _ string, _ decimal.Decimal) error {
	return errors.New("mock: UpdateTickSize not used")
}
func (m *mockMarketRepo) ListActive(_ context.Context) ([]*market.Market, error) {
	return nil, errMockNotFound
}

func newBTCMarket(t *testing.T, windowStart time.Time) *market.Market {
	t.Helper()
	m, err := market.New(market.Params{
		ID:          "evt-btc",
		Asset:       market.BTC,
		WindowStart: windowStart,
		ConditionID: polyid.ConditionID("0xcondbtc"),
		UpTokenID:   polyid.TokenID("0xupbtc"),
		DownTokenID: polyid.TokenID("0xdownbtc"),
		TickSize:    decimal.NewFromFloat(0.01),
		FeeEnabled:  true,
	})
	require.NoError(t, err)
	return m
}

func TestGetActiveMarket_Execute(t *testing.T) {
	t.Parallel()

	windowStart := time.Now().UTC().Truncate(5 * time.Minute)

	tests := []struct {
		name      string
		input     dto.Input
		repo      mockMarketRepo
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output)
	}{
		{
			name:  "returns active BTC market",
			input: dto.Input{Asset: "btc", WindowStart: windowStart},
			repo:  mockMarketRepo{found: newBTCMarket(t, windowStart)},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "btc", out.Asset)
				assert.Equal(t, "0xcondbtc", out.ConditionID)
				assert.NotEmpty(t, out.WindowStart)
				assert.NotEmpty(t, out.WindowEnd)
			},
		},
		{
			name:      "empty asset returns client error",
			input:     dto.Input{Asset: "", WindowStart: windowStart},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "market not found returns not found error",
			input:     dto.Input{Asset: "btc", WindowStart: windowStart},
			repo:      mockMarketRepo{foundErr: errMockNotFound},
			wantErr:   true,
			errTarget: &errtypes.NotFoundError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := getactivemarket.New(&tt.repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				if tt.checkOut != nil {
					tt.checkOut(t, out)
				}
			}
		})
	}
}
