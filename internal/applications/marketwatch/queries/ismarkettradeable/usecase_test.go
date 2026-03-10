package ismarkettradeable_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/ismarkettradeable"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/ismarkettradeable/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// errMockNotFound satisfies nilnil linter for unused mock methods.
var errMockNotFound = errors.New("mock: not found")

type mockMarketRepo struct {
	activeMarkets []*market.Market
	listErr       error
}

func (m *mockMarketRepo) Save(_ context.Context, _ *market.Market) error {
	return errors.New("mock: Save not used")
}
func (m *mockMarketRepo) FindByAssetAndWindow(_ context.Context, _ market.Asset, _ time.Time) (*market.Market, error) {
	return nil, errMockNotFound
}
func (m *mockMarketRepo) UpdateTickSize(_ context.Context, _ string, _ decimal.Decimal) error {
	return errors.New("mock: UpdateTickSize not used")
}
func (m *mockMarketRepo) ListActive(_ context.Context) ([]*market.Market, error) {
	return m.activeMarkets, m.listErr
}

func newMarketWithCond(t *testing.T, asset market.Asset, conditionID string, feeEnabled bool) *market.Market {
	t.Helper()
	m, err := market.New(market.Params{
		ID:          "evt-" + conditionID,
		Asset:       asset,
		WindowStart: time.Now().UTC().Truncate(5 * time.Minute),
		ConditionID: polyid.ConditionID(conditionID),
		UpTokenID:   polyid.TokenID("0xup"),
		DownTokenID: polyid.TokenID("0xdown"),
		TickSize:    decimal.NewFromFloat(0.01),
		FeeEnabled:  feeEnabled,
	})
	require.NoError(t, err)
	return m
}

func TestIsMarketTradeable_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         dto.Input
		activeMarkets []*market.Market
		listErr       error
		wantErr       bool
		errTarget     any
		wantTradeable bool
		wantReason    string
	}{
		{
			name:  "active market with fees enabled is tradeable",
			input: dto.Input{ConditionID: "0xcond1"},
			activeMarkets: []*market.Market{
				newMarketWithCond(t, market.BTC, "0xcond1", true),
			},
			wantTradeable: true,
			wantReason:    "",
		},
		{
			name:  "market not in active list is not tradeable",
			input: dto.Input{ConditionID: "0xcond_missing"},
			activeMarkets: []*market.Market{
				newMarketWithCond(t, market.BTC, "0xcond1", true),
			},
			wantTradeable: false,
			wantReason:    "market not active",
		},
		{
			name:  "market with fees disabled is not tradeable",
			input: dto.Input{ConditionID: "0xcond_nofee"},
			activeMarkets: []*market.Market{
				newMarketWithCond(t, market.ETH, "0xcond_nofee", false),
			},
			wantTradeable: false,
			wantReason:    "fees not enabled",
		},
		{
			name:      "empty condition ID returns client error",
			input:     dto.Input{ConditionID: ""},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "list active error returns internal error",
			input:     dto.Input{ConditionID: "0xcond1"},
			listErr:   errors.New("db unavailable"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockMarketRepo{
				activeMarkets: tt.activeMarkets,
				listErr:       tt.listErr,
			}
			uc := ismarkettradeable.New(repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantTradeable, out.Tradeable)
				assert.Equal(t, tt.wantReason, out.Reason)
			}
		})
	}
}
