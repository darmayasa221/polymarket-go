package startwindow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// mockWindowStateStore implements ports.WindowStateStore.
type mockWindowStateStore struct {
	saved   windowstate.WindowState
	saveErr error
	getErr  error
}

func (m *mockWindowStateStore) SaveWindowState(_ context.Context, state windowstate.WindowState) error {
	m.saved = state
	return m.saveErr
}
func (m *mockWindowStateStore) GetWindowState(_ context.Context, _ string) (windowstate.WindowState, error) {
	return windowstate.WindowState{}, m.getErr
}

func TestStartWindow_Execute(t *testing.T) {
	t.Parallel()

	validInput := dto.Input{
		Asset:       "btc",
		MarketID:    "evt-btc-123",
		ConditionID: "0xcond",
		UpTokenID:   "0xup",
		DownTokenID: "0xdown",
		TickSize:    "0.01",
		OpenPrice:   decimal.NewFromFloat(67000),
	}

	tests := []struct {
		name      string
		input     dto.Input
		saveErr   error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output, store *mockWindowStateStore)
	}{
		{
			name:  "valid input initializes window state",
			input: validInput,
			checkOut: func(t *testing.T, out dto.Output, store *mockWindowStateStore) {
				t.Helper()
				assert.Equal(t, "btc", out.Asset)
				assert.Equal(t, "evt-btc-123", out.MarketID)
				assert.Equal(t, "0xcond", out.ConditionID)
				assert.False(t, out.WindowStart.IsZero())
				assert.False(t, out.WindowEnd.IsZero())
				assert.True(t, out.WindowEnd.After(out.WindowStart))
				// WindowEnd should be exactly 5 minutes after WindowStart
				assert.Equal(t, 5*60, int(out.WindowEnd.Sub(out.WindowStart).Seconds()))
				// State was persisted
				assert.Equal(t, windowstate.WindowOpen, store.saved.Status)
				assert.Equal(t, "btc", store.saved.Asset)
			},
		},
		{
			name:      "empty asset returns client error",
			input:     dto.Input{Asset: "", MarketID: "evt-123", ConditionID: "0x", OpenPrice: decimal.NewFromFloat(1)},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "unsupported asset returns client error",
			input:     dto.Input{Asset: "doge", MarketID: "evt-123", ConditionID: "0x", OpenPrice: decimal.NewFromFloat(1)},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "empty marketID returns client error",
			input:     dto.Input{Asset: "btc", MarketID: "", ConditionID: "0x", OpenPrice: decimal.NewFromFloat(1)},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "store save failure returns internal error",
			input:     validInput,
			saveErr:   errors.New("store unavailable"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := &mockWindowStateStore{saveErr: tt.saveErr}
			uc := startwindow.New(store)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				if tt.checkOut != nil {
					tt.checkOut(t, out, store)
				}
			}
		})
	}
}
