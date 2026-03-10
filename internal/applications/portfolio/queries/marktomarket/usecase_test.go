package marktomarket_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/marktomarket"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/marktomarket/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

var errMockNotFound = errors.New("mock: not found")

type mockPositionRepo struct {
	positions []*position.Position
	listErr   error
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
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.positions, nil
}
func (m *mockPositionRepo) ListClosed(_ context.Context) ([]*position.Position, error) {
	return nil, errMockNotFound
}
func (m *mockPositionRepo) ListClosedWithExitPrice(_ context.Context) ([]portfolioports.ClosedPositionRecord, error) {
	return nil, errMockNotFound
}
func (m *mockPositionRepo) Close(_ context.Context, _ string, _ decimal.Decimal, _ time.Time) error {
	return errors.New("mock: Close not used")
}

func makePosition(t *testing.T, tokenID, avgPrice string) *position.Position {
	t.Helper()
	pos, err := position.New(position.Params{
		Asset:    market.Asset("btc"),
		TokenID:  polyid.TokenID(tokenID),
		Outcome:  market.Up,
		Size:     decimal.RequireFromString("10"),
		AvgPrice: decimal.RequireFromString(avgPrice),
		MarketID: "market-btc-1",
	})
	require.NoError(t, err)
	return pos
}

func TestMarkToMarket_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     dto.Input
		positions []*position.Position
		listErr   error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output)
	}{
		{
			name:      "ListOpen failure returns internal error",
			input:     dto.Input{Prices: map[string]string{}},
			listErr:   errors.New("db down"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:      "no open positions returns zero total",
			input:     dto.Input{Prices: map[string]string{}},
			positions: []*position.Position{},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Empty(t, out.Marks)
				assert.Equal(t, "0", out.TotalUnrealisedPnL)
			},
		},
		{
			name: "position with current price above avg gives positive PnL",
			input: dto.Input{Prices: map[string]string{
				"token-up": "0.80",
			}},
			positions: []*position.Position{makePosition(t, "token-up", "0.60")},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				require.Len(t, out.Marks, 1)
				// size=10, avgPrice=0.60, currentPrice=0.80 → PnL=10*(0.80-0.60)=2.00
				pnl := decimal.RequireFromString(out.Marks[0].UnrealisedPnL)
				assert.True(t, pnl.IsPositive())
				assert.Equal(t, "2", out.TotalUnrealisedPnL)
			},
		},
		{
			name:      "missing price falls back to avg price giving zero PnL",
			input:     dto.Input{Prices: map[string]string{}},
			positions: []*position.Position{makePosition(t, "token-up", "0.60")},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				require.Len(t, out.Marks, 1)
				assert.Equal(t, "0", out.Marks[0].UnrealisedPnL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockPositionRepo{positions: tt.positions, listErr: tt.listErr}
			uc := marktomarket.New(repo)

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
