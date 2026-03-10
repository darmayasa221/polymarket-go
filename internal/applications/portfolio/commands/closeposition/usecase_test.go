package closeposition_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/closeposition"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/closeposition/dto"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

var errMockNotFound = errors.New("mock: not found")

type mockPositionRepo struct {
	pos      *position.Position
	findErr  error
	closeErr error
	closed   string
}

func (m *mockPositionRepo) Save(_ context.Context, _ *position.Position) error {
	return errors.New("mock: Save not used")
}
func (m *mockPositionRepo) FindByID(_ context.Context, _ string) (*position.Position, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.pos, nil
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
	return nil, errMockNotFound
}
func (m *mockPositionRepo) Close(_ context.Context, posID string, _ decimal.Decimal, _ time.Time) error {
	if m.closeErr != nil {
		return m.closeErr
	}
	m.closed = posID
	return nil
}

func makePosition(t *testing.T) *position.Position {
	t.Helper()
	pos, err := position.New(position.Params{
		Asset:    market.Asset("btc"),
		TokenID:  polyid.TokenID("token-up"),
		Outcome:  market.Up,
		Size:     decimal.RequireFromString("10"),
		AvgPrice: decimal.RequireFromString("0.60"),
		MarketID: "market-btc-1",
	})
	require.NoError(t, err)
	return pos
}

func TestClosePosition_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     dto.Input
		pos       *position.Position
		findErr   error
		closeErr  error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output)
	}{
		{
			name:      "empty position ID returns client error",
			input:     dto.Input{PositionID: "", ExitPrice: "0.80"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "zero exit price returns client error",
			input:     dto.Input{PositionID: "pos-1", ExitPrice: "0"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "position not found returns client error",
			input:     dto.Input{PositionID: "pos-1", ExitPrice: "0.80"},
			findErr:   errMockNotFound,
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:  "stop loss exit returns negative realized PnL",
			input: dto.Input{PositionID: "pos-1", ExitPrice: "0.40"},
			pos:   makePosition(t),
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "pos-1", out.PositionID)
				// size=10, avgPrice=0.60, exitPrice=0.40 → PnL = 10*(0.40-0.60) = -2.00
				pnl := decimal.RequireFromString(out.RealisedPnL)
				assert.True(t, pnl.IsNegative(), "stop loss should be negative PnL")
			},
		},
		{
			name:  "take profit exit returns positive realized PnL",
			input: dto.Input{PositionID: "pos-1", ExitPrice: "0.80"},
			pos:   makePosition(t),
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				// size=10, avgPrice=0.60, exitPrice=0.80 → PnL = 10*(0.80-0.60) = 2.00
				pnl := decimal.RequireFromString(out.RealisedPnL)
				assert.True(t, pnl.IsPositive(), "take profit should be positive PnL")
				assert.Equal(t, "2", pnl.String())
			},
		},
		{
			name:      "repo Close failure returns internal error",
			input:     dto.Input{PositionID: "pos-1", ExitPrice: "0.80"},
			pos:       makePosition(t),
			closeErr:  errors.New("db down"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockPositionRepo{pos: tt.pos, findErr: tt.findErr, closeErr: tt.closeErr}
			uc := closeposition.New(repo)

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
