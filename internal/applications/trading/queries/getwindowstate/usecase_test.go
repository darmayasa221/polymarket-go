package getwindowstate_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

type mockWindowStore struct {
	state  windowstate.WindowState
	getErr error
}

func (m *mockWindowStore) GetWindowState(_ context.Context, _ string) (windowstate.WindowState, error) {
	if m.getErr != nil {
		return windowstate.WindowState{}, m.getErr
	}
	return m.state, nil
}

func (m *mockWindowStore) SaveWindowState(_ context.Context, _ windowstate.WindowState) error {
	return errors.New("mock: SaveWindowState not used")
}

func TestGetWindowState_Execute(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	full := windowstate.WindowState{
		Asset:       "eth",
		MarketID:    "market-eth",
		ConditionID: "cond-eth",
		WindowStart: now,
		WindowEnd:   now.Add(5 * time.Minute),
		Status:      windowstate.WindowOpen,
		OpenOrders:  []windowstate.OrderSummary{},
	}

	tests := []struct {
		name      string
		input     dto.Input
		state     windowstate.WindowState
		getErr    error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output)
	}{
		{
			name:      "empty asset returns client error",
			input:     dto.Input{Asset: ""},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "store error returns internal error",
			input:     dto.Input{Asset: "eth"},
			getErr:    errors.New("store unavailable"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:  "valid asset returns window state",
			input: dto.Input{Asset: "eth"},
			state: full,
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "eth", out.Asset)
				assert.Equal(t, "market-eth", out.MarketID)
				assert.Equal(t, windowstate.WindowOpen, out.Status)
				assert.Equal(t, now, out.WindowStart)
				assert.Equal(t, now.Add(5*time.Minute), out.WindowEnd)
			},
		},
		{
			name:  "closed window state returned correctly",
			input: dto.Input{Asset: "btc"},
			state: windowstate.WindowState{
				Asset:  "btc",
				Status: windowstate.WindowClosed,
			},
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, windowstate.WindowClosed, out.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := &mockWindowStore{state: tt.state, getErr: tt.getErr}
			uc := getwindowstate.New(store)

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
