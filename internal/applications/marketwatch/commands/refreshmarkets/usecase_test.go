package refreshmarkets_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/refreshmarkets"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/refreshmarkets/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// errMockNotFound satisfies nilnil linter for unused mock methods.
var errMockNotFound = errors.New("mock: not found")

// mockMarketSource implements ports.MarketSource.
type mockMarketSource struct {
	markets []*market.Market
	err     error
}

func (m *mockMarketSource) FetchActive5mMarkets(_ context.Context) ([]*market.Market, error) {
	return m.markets, m.err
}

// mockMarketRepo implements ports.MarketRepository.
type mockMarketRepo struct {
	saved   []*market.Market
	saveErr error
}

func (m *mockMarketRepo) Save(_ context.Context, mk *market.Market) error {
	m.saved = append(m.saved, mk)
	return m.saveErr
}
func (m *mockMarketRepo) FindByAssetAndWindow(_ context.Context, _ market.Asset, _ time.Time) (*market.Market, error) {
	return nil, errMockNotFound
}
func (m *mockMarketRepo) UpdateTickSize(_ context.Context, _ string, _ decimal.Decimal) error {
	return errors.New("mock: UpdateTickSize not used in RefreshMarkets")
}
func (m *mockMarketRepo) ListActive(_ context.Context) ([]*market.Market, error) {
	return nil, errMockNotFound
}

func newTestMarket(t *testing.T, asset market.Asset) *market.Market {
	t.Helper()
	m, err := market.New(market.Params{
		ID:          "evt-" + string(asset),
		Asset:       asset,
		WindowStart: time.Now().UTC().Truncate(5 * time.Minute),
		ConditionID: polyid.ConditionID("0xabc" + string(asset)),
		UpTokenID:   polyid.TokenID("0xup" + string(asset)),
		DownTokenID: polyid.TokenID("0xdown" + string(asset)),
		TickSize:    decimal.NewFromFloat(0.01),
		FeeEnabled:  true,
	})
	require.NoError(t, err)
	return m
}

func TestRefreshMarkets_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		sourceMarkets []*market.Market
		sourceErr     error
		saveErr       error
		wantErr       bool
		errTarget     any
		wantRefreshed int
	}{
		{
			name: "fetches and saves 4 active markets",
			sourceMarkets: []*market.Market{
				newTestMarket(t, market.BTC),
				newTestMarket(t, market.ETH),
				newTestMarket(t, market.SOL),
				newTestMarket(t, market.XRP),
			},
			wantRefreshed: 4,
		},
		{
			name:          "fetches 1 market successfully",
			sourceMarkets: []*market.Market{newTestMarket(t, market.BTC)},
			wantRefreshed: 1,
		},
		{
			name:          "empty market list returns zero refreshed",
			sourceMarkets: []*market.Market{},
			wantRefreshed: 0,
		},
		{
			name:      "source fetch failure returns internal error",
			sourceErr: errors.New("gamma api unavailable"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name: "save failure returns internal error",
			sourceMarkets: []*market.Market{
				newTestMarket(t, market.BTC),
			},
			saveErr:   errors.New("db write failed"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			source := &mockMarketSource{markets: tt.sourceMarkets, err: tt.sourceErr}
			repo := &mockMarketRepo{saveErr: tt.saveErr}
			uc := refreshmarkets.New(source, repo)

			out, err := uc.Execute(t.Context(), dto.Input{})

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantRefreshed, out.Refreshed)
				assert.Len(t, out.Assets, tt.wantRefreshed)
				assert.Equal(t, tt.wantRefreshed, len(repo.saved))
			}
		})
	}
}
