package updateticksize_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/updateticksize"
	"github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/updateticksize/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// errMockNotFound satisfies nilnil linter for unused mock methods.
var errMockNotFound = errors.New("mock: not found")

type mockMarketRepo struct {
	updateErr error
	updated   string
	newTick   decimal.Decimal
}

func (m *mockMarketRepo) Save(_ context.Context, _ *market.Market) error {
	return errors.New("mock: Save not used in UpdateTickSize")
}
func (m *mockMarketRepo) FindByAssetAndWindow(_ context.Context, _ market.Asset, _ time.Time) (*market.Market, error) {
	return nil, errMockNotFound
}
func (m *mockMarketRepo) UpdateTickSize(_ context.Context, conditionID string, newTick decimal.Decimal) error {
	m.updated = conditionID
	m.newTick = newTick
	return m.updateErr
}
func (m *mockMarketRepo) ListActive(_ context.Context) ([]*market.Market, error) {
	return nil, errMockNotFound
}

func TestUpdateTickSize_Execute(t *testing.T) {
	t.Parallel()

	validInput := dto.Input{
		ConditionID: "0xabc123",
		TokenID:     "0xtoken",
		OldTickSize: "0.01",
		NewTickSize: "0.1",
	}

	tests := []struct {
		name      string
		input     dto.Input
		updateErr error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output, repo *mockMarketRepo)
	}{
		{
			name:  "valid tick size 0.1 is accepted",
			input: validInput,
			checkOut: func(t *testing.T, out dto.Output, repo *mockMarketRepo) {
				t.Helper()
				assert.Equal(t, "0xabc123", out.ConditionID)
				assert.Equal(t, "0.1", out.NewTickSize)
				assert.Equal(t, "0xabc123", repo.updated)
			},
		},
		{
			name:  "valid tick size 0.01 is accepted",
			input: dto.Input{ConditionID: "0xabc", NewTickSize: "0.01"},
			checkOut: func(t *testing.T, _ dto.Output, _ *mockMarketRepo) {
				t.Helper()
			},
		},
		{
			name:  "valid tick size 0.001 is accepted",
			input: dto.Input{ConditionID: "0xabc", NewTickSize: "0.001"},
			checkOut: func(t *testing.T, _ dto.Output, _ *mockMarketRepo) {
				t.Helper()
			},
		},
		{
			name:  "valid tick size 0.0001 is accepted",
			input: dto.Input{ConditionID: "0xabc", NewTickSize: "0.0001"},
			checkOut: func(t *testing.T, _ dto.Output, _ *mockMarketRepo) {
				t.Helper()
			},
		},
		{
			name:      "empty condition ID returns client error",
			input:     dto.Input{ConditionID: "", NewTickSize: "0.01"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "invalid tick size 0.05 returns client error",
			input:     dto.Input{ConditionID: "0xabc", NewTickSize: "0.05"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "non-numeric tick size returns client error",
			input:     dto.Input{ConditionID: "0xabc", NewTickSize: "bad"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "repo update error returns internal error",
			input:     validInput,
			updateErr: errors.New("db write failed"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockMarketRepo{updateErr: tt.updateErr}
			uc := updateticksize.New(repo)

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
