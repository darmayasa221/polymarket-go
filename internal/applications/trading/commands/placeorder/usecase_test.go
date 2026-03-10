package placeorder_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// errMockNotFound satisfies nilnil linter.
var errMockNotFound = errors.New("mock: not found")

// mockOrderRepo implements ports.OrderRepository.
type mockOrderRepo struct {
	saved   *order.Order
	saveErr error
}

func (m *mockOrderRepo) Save(_ context.Context, o *order.Order) error {
	m.saved = o
	return m.saveErr
}
func (m *mockOrderRepo) FindByID(_ context.Context, _ polyid.OrderID) (*order.Order, error) {
	return nil, errMockNotFound
}
func (m *mockOrderRepo) ListOpenByMarket(_ context.Context, _ string) ([]*order.Order, error) {
	return nil, errMockNotFound
}
func (m *mockOrderRepo) UpdateStatus(_ context.Context, _ polyid.OrderID, _ order.OrderStatus) error {
	return errors.New("mock: UpdateStatus not used")
}

// mockWindowStateStore implements ports.WindowStateStore.
type mockWindowStateStore struct {
	state  windowstate.WindowState
	getErr error
}

func (m *mockWindowStateStore) SaveWindowState(_ context.Context, _ windowstate.WindowState) error {
	return errors.New("mock: SaveWindowState not used in PlaceOrder")
}
func (m *mockWindowStateStore) GetWindowState(_ context.Context, _ string) (windowstate.WindowState, error) {
	if m.getErr != nil {
		return windowstate.WindowState{}, m.getErr
	}
	return m.state, nil
}

func openWindowState() windowstate.WindowState {
	now := timeutil.Now()
	return windowstate.WindowState{
		MarketID:    "evt-btc",
		Asset:       "btc",
		WindowStart: timeutil.WindowStart(now),
		WindowEnd:   timeutil.WindowEnd(now),
		ConditionID: "0xcond",
		UpTokenID:   "12345",
		DownTokenID: "67890",
		TickSize:    "0.01",
		OpenPrice:   decimal.NewFromFloat(67000),
		Status:      windowstate.WindowOpen,
		OpenOrders:  []windowstate.OrderSummary{},
	}
}

func TestPlaceOrder_Execute(t *testing.T) {
	t.Parallel()

	validInput := dto.Input{
		Asset:         "btc",
		Outcome:       "Up",
		Side:          "buy",
		Price:         decimal.NewFromFloat(0.55),
		Size:          decimal.NewFromFloat(10),
		TokenID:       "12345",
		FeeRateBps:    156,
		FunderAddress: "0xfunder",
	}

	tests := []struct {
		name      string
		input     dto.Input
		state     windowstate.WindowState
		stateErr  error
		saveErr   error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output, repo *mockOrderRepo)
	}{
		{
			name:  "valid buy order returns unsigned hash",
			input: validInput,
			state: openWindowState(),
			checkOut: func(t *testing.T, out dto.Output, repo *mockOrderRepo) {
				t.Helper()
				assert.NotEmpty(t, out.OrderID)
				assert.Len(t, out.UnsignedHash, 32, "EIP-712 hash must be 32 bytes")
				assert.NotZero(t, out.GTDExpiry)
				assert.NotNil(t, repo.saved)
			},
		},
		{
			name:  "valid sell order returns unsigned hash",
			input: dto.Input{Asset: "eth", Outcome: "Down", Side: "sell", Price: decimal.NewFromFloat(0.40), Size: decimal.NewFromFloat(5), TokenID: "99", FeeRateBps: 100, FunderAddress: "0xf"},
			state: windowstate.WindowState{Status: windowstate.WindowOpen, MarketID: "evt-eth", Asset: "eth", UpTokenID: "88", DownTokenID: "99", WindowEnd: timeutil.WindowEnd(timeutil.Now())},
			checkOut: func(t *testing.T, out dto.Output, _ *mockOrderRepo) {
				t.Helper()
				assert.Len(t, out.UnsignedHash, 32)
			},
		},
		{
			name:      "window not open returns client error",
			input:     validInput,
			state:     windowstate.WindowState{Status: windowstate.WindowClosed},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "state not found returns internal error",
			input:     validInput,
			stateErr:  errMockNotFound,
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:      "invalid outcome returns client error",
			input:     dto.Input{Asset: "btc", Outcome: "Yes", Side: "buy", Price: decimal.NewFromFloat(0.50), Size: decimal.NewFromFloat(10), TokenID: "12345", FeeRateBps: 100, FunderAddress: "0xf"},
			state:     openWindowState(),
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "invalid side returns client error",
			input:     dto.Input{Asset: "btc", Outcome: "Up", Side: "hold", Price: decimal.NewFromFloat(0.50), Size: decimal.NewFromFloat(10), TokenID: "12345", FeeRateBps: 100, FunderAddress: "0xf"},
			state:     openWindowState(),
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "repo save failure returns internal error",
			input:     validInput,
			state:     openWindowState(),
			saveErr:   errors.New("db write failed"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockOrderRepo{saveErr: tt.saveErr}
			store := &mockWindowStateStore{state: tt.state, getErr: tt.stateErr}
			uc := placeorder.New(repo, store)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
				assert.Empty(t, out.OrderID)
			} else {
				require.NoError(t, err)
				if tt.checkOut != nil {
					tt.checkOut(t, out, repo)
				}
			}
		})
	}
}
