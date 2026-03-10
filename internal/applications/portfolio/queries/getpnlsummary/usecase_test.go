package getpnlsummary_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpnlsummary"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpnlsummary/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

var errMockNotFound = errors.New("mock: not found")

type mockPositionRepo struct {
	records []portfolioports.ClosedPositionRecord
	listErr error
}

func (m *mockPositionRepo) Save(_ context.Context, _ *position.Position) error {
	return errors.New("mock: Save not used")
}
func (m *mockPositionRepo) FindByID(_ context.Context, _ string) (*position.Position, error) {
	return nil, errMockNotFound
}
func (m *mockPositionRepo) FindByMarket(_ context.Context, _ string) ([]*position.Position, error) {
	return nil, errMockNotFound
}
func (m *mockPositionRepo) ListOpen(_ context.Context) ([]*position.Position, error) {
	return nil, errMockNotFound
}
func (m *mockPositionRepo) ListClosed(_ context.Context) ([]*position.Position, error) {
	return nil, errMockNotFound
}
func (m *mockPositionRepo) ListClosedWithExitPrice(_ context.Context) ([]portfolioports.ClosedPositionRecord, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.records, nil
}
func (m *mockPositionRepo) Close(_ context.Context, _ string, _ decimal.Decimal, _ time.Time) error {
	return errors.New("mock: Close not used")
}

func makeClosedRecord(t *testing.T, avgPrice, exitPrice string) portfolioports.ClosedPositionRecord {
	t.Helper()
	pos, err := position.New(position.Params{
		Asset:    market.Asset("btc"),
		TokenID:  polyid.TokenID("token-up"),
		Outcome:  market.Up,
		Size:     decimal.RequireFromString("10"),
		AvgPrice: decimal.RequireFromString(avgPrice),
		MarketID: "market-btc-1",
	})
	require.NoError(t, err)
	return portfolioports.ClosedPositionRecord{
		Pos:       pos,
		ExitPrice: decimal.RequireFromString(exitPrice),
		ClosedAt:  time.Now().UTC(),
	}
}

func TestGetPnLSummary_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		records   []portfolioports.ClosedPositionRecord
		listErr   error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output)
	}{
		{
			name:      "ListClosedWithExitPrice failure returns internal error",
			listErr:   errors.New("db down"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:    "no closed positions returns zero summary",
			records: []portfolioports.ClosedPositionRecord{},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "0", out.TotalRealisedPnL)
				assert.Equal(t, "0", out.TotalUnrealisedPnL)
				assert.Equal(t, 0, out.WinCount)
				assert.Equal(t, 0, out.LossCount)
				assert.Equal(t, 0, out.TotalCount)
			},
		},
		{
			name: "one win one loss returns correct totals",
			records: []portfolioports.ClosedPositionRecord{
				makeClosedRecord(t, "0.60", "0.80"),
				makeClosedRecord(t, "0.60", "0.40"),
			},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "0", out.TotalRealisedPnL)
				assert.Equal(t, "0", out.TotalUnrealisedPnL)
				assert.Equal(t, 1, out.WinCount)
				assert.Equal(t, 1, out.LossCount)
				assert.Equal(t, 2, out.TotalCount)
			},
		},
		{
			name: "two wins sums correctly",
			records: []portfolioports.ClosedPositionRecord{
				makeClosedRecord(t, "0.60", "0.80"), // +2
				makeClosedRecord(t, "0.50", "0.75"), // +2.5
			},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				total := decimal.RequireFromString(out.TotalRealisedPnL)
				assert.True(t, total.IsPositive())
				assert.Equal(t, 2, out.WinCount)
				assert.Equal(t, 0, out.LossCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockPositionRepo{records: tt.records, listErr: tt.listErr}
			uc := getpnlsummary.New(repo)

			out, err := uc.Execute(t.Context(), dto.Input{})

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
