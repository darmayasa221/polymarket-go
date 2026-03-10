package getpositions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpositions"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpositions/dto"
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

func makePosition(t *testing.T, asset market.Asset) *position.Position {
	t.Helper()
	pos, err := position.New(position.Params{
		Asset:    asset,
		TokenID:  polyid.TokenID("token-" + string(asset)),
		Outcome:  market.Up,
		Size:     decimal.RequireFromString("10"),
		AvgPrice: decimal.RequireFromString("0.60"),
		MarketID: "market-" + string(asset),
	})
	require.NoError(t, err)
	return pos
}

func TestGetPositions_Execute(t *testing.T) {
	t.Parallel()

	btcPos := makePosition(t, "btc")
	ethPos := makePosition(t, "eth")

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
			input:     dto.Input{},
			listErr:   errors.New("db down"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:      "no filter returns all positions",
			input:     dto.Input{},
			positions: []*position.Position{btcPos, ethPos},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Len(t, out.Positions, 2)
			},
		},
		{
			name:      "asset filter returns only matching positions",
			input:     dto.Input{Asset: "btc"},
			positions: []*position.Position{btcPos, ethPos},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				require.Len(t, out.Positions, 1)
				assert.Equal(t, "btc", out.Positions[0].Asset)
			},
		},
		{
			name:      "filter with no matches returns empty slice",
			input:     dto.Input{Asset: "sol"},
			positions: []*position.Position{btcPos, ethPos},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Empty(t, out.Positions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockPositionRepo{positions: tt.positions, listErr: tt.listErr}
			uc := getpositions.New(repo)

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
