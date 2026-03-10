package getcurrentsignal_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal"
	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// errMockNotFound is returned by unused mock methods to satisfy the nilnil linter.
var errMockNotFound = errors.New("mock: not found")

// mockPriceRepo implements ports.PriceRepository for testing.
type mockPriceRepo struct {
	latestChainlink *oracle.Price
	latestChainErr  error
	latestAny       *oracle.Price
	latestAnyErr    error
	windowOpen      *oracle.Price
	windowOpenErr   error
}

func (m *mockPriceRepo) Save(_ context.Context, _ *oracle.Price) error {
	return errors.New("mock: Save not used in query")
}
func (m *mockPriceRepo) LatestByAsset(_ context.Context, _ string) (*oracle.Price, error) {
	if m.latestAny != nil {
		return m.latestAny, m.latestAnyErr
	}
	return nil, errMockNotFound
}
func (m *mockPriceRepo) LatestChainlinkByAsset(_ context.Context, _ string) (*oracle.Price, error) {
	if m.latestChainErr != nil {
		return nil, m.latestChainErr
	}
	if m.latestChainlink != nil {
		return m.latestChainlink, nil
	}
	return nil, errMockNotFound
}
func (m *mockPriceRepo) WindowOpenPrice(_ context.Context, _ string, _ time.Time) (*oracle.Price, error) {
	if m.windowOpenErr != nil {
		return nil, m.windowOpenErr
	}
	if m.windowOpen != nil {
		return m.windowOpen, nil
	}
	return nil, errMockNotFound
}

func newPrice(t *testing.T, asset, source, value string) *oracle.Price {
	t.Helper()
	v, err := decimal.NewFromString(value)
	require.NoError(t, err)
	src := oracle.PriceSource(source)
	p, err := oracle.New(oracle.Params{
		Asset:      asset,
		Source:     src,
		Value:      v,
		ReceivedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	return p
}

func TestGetCurrentSignal_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         dto.Input
		repo          mockPriceRepo
		wantErr       bool
		errTarget     any
		wantPredicted string
		wantConfMin   float64 // minimum expected confidence
	}{
		{
			name:  "5% up move returns Up with full confidence",
			input: dto.Input{Asset: "btc"},
			repo: mockPriceRepo{
				windowOpen:      newPrice(t, "btc", "chainlink", "100000"),
				latestChainlink: newPrice(t, "btc", "chainlink", "105000"), // +5%
			},
			wantPredicted: "Up",
			wantConfMin:   1.0, // exactly 5% = capped at 1.0
		},
		{
			name:  "2.5% down move returns Down with 0.5 confidence",
			input: dto.Input{Asset: "btc"},
			repo: mockPriceRepo{
				windowOpen:      newPrice(t, "btc", "chainlink", "100000"),
				latestChainlink: newPrice(t, "btc", "chainlink", "97500"), // -2.5%
			},
			wantPredicted: "Down",
			wantConfMin:   0.49,
		},
		{
			name:  "flat price returns confidence near zero",
			input: dto.Input{Asset: "eth"},
			repo: mockPriceRepo{
				windowOpen:      newPrice(t, "eth", "chainlink", "3000"),
				latestChainlink: newPrice(t, "eth", "chainlink", "3000"),
			},
			wantPredicted: "Up", // equal price → Up (GreaterThanOrEqual)
			wantConfMin:   0.0,
		},
		{
			name:      "empty asset returns client error",
			input:     dto.Input{Asset: ""},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:  "missing window open price returns internal error",
			input: dto.Input{Asset: "btc"},
			repo: mockPriceRepo{
				windowOpenErr: errMockNotFound,
			},
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:  "chainlink unavailable falls back to binance",
			input: dto.Input{Asset: "sol"},
			repo: mockPriceRepo{
				windowOpen:   newPrice(t, "sol", "chainlink", "200"),
				latestAny:    newPrice(t, "sol", "binance", "210"),
				latestAnyErr: nil,
			},
			wantPredicted: "Up",
			wantConfMin:   1.0, // 5% move
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := getcurrentsignal.New(&tt.repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPredicted, out.Signal.Predicted)
				f, _ := out.Signal.Confidence.Float64()
				assert.GreaterOrEqual(t, f, tt.wantConfMin)
				assert.LessOrEqual(t, f, 1.0)
			}
		})
	}
}
