package openposition_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/openposition"
	"github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/openposition/dto"
	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

var errMockNotFound = errors.New("mock: not found")

type mockPositionRepo struct {
	saveErr error
	saved   *position.Position
}

func (m *mockPositionRepo) Save(_ context.Context, pos *position.Position) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = pos
	return nil
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
	return nil, errMockNotFound
}
func (m *mockPositionRepo) Close(_ context.Context, _ string, _ decimal.Decimal, _ time.Time) error {
	return errors.New("mock: Close not used")
}

func TestOpenPosition_Execute(t *testing.T) {
	t.Parallel()

	validInput := dto.Input{
		Asset:    "btc",
		TokenID:  "token-up-btc",
		Outcome:  "Up",
		Size:     "15",
		AvgPrice: "0.62",
		MarketID: "market-btc-1",
	}

	tests := []struct {
		name      string
		input     dto.Input
		saveErr   error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output, repo *mockPositionRepo)
	}{
		{
			name:  "valid input creates position",
			input: validInput,
			checkOut: func(t *testing.T, out dto.Output, repo *mockPositionRepo) {
				t.Helper()
				assert.NotEmpty(t, out.PositionID)
				require.NotNil(t, repo.saved)
			},
		},
		{
			name:      "invalid asset returns client error",
			input:     dto.Input{Asset: "invalid", TokenID: "tok", Outcome: "Up", Size: "10", AvgPrice: "0.5", MarketID: "m1"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "invalid outcome returns client error",
			input:     dto.Input{Asset: "btc", TokenID: "tok", Outcome: "Yes", Size: "10", AvgPrice: "0.5", MarketID: "m1"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "zero size returns client error",
			input:     dto.Input{Asset: "btc", TokenID: "tok", Outcome: "Up", Size: "0", AvgPrice: "0.5", MarketID: "m1"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "zero avg price returns client error",
			input:     dto.Input{Asset: "btc", TokenID: "tok", Outcome: "Up", Size: "10", AvgPrice: "0", MarketID: "m1"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "repo save failure returns internal error",
			input:     validInput,
			saveErr:   errors.New("db down"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockPositionRepo{saveErr: tt.saveErr}
			uc := openposition.New(repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				if tt.checkOut != nil {
					tt.checkOut(t, out, repo)
				}
			}
		})
	}
}
