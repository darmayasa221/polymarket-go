package closewindow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// errMockNotFound satisfies nilnil linter.
var errMockNotFound = errors.New("mock: not found")

// --- mocks ---

type mockWindowStore struct {
	state   windowstate.WindowState
	getErr  error
	saveErr error
	saved   *windowstate.WindowState
}

func (m *mockWindowStore) GetWindowState(_ context.Context, _ string) (windowstate.WindowState, error) {
	if m.getErr != nil {
		return windowstate.WindowState{}, m.getErr
	}
	return m.state, nil
}

func (m *mockWindowStore) SaveWindowState(_ context.Context, s windowstate.WindowState) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = &s
	return nil
}

type mockOrderRepo struct {
	orders  []*order.Order
	listErr error
	updated []polyid.OrderID
}

func (m *mockOrderRepo) Save(_ context.Context, _ *order.Order) error {
	return errors.New("mock: Save not used")
}
func (m *mockOrderRepo) FindByID(_ context.Context, _ polyid.OrderID) (*order.Order, error) {
	return nil, errMockNotFound
}
func (m *mockOrderRepo) ListOpenByMarket(_ context.Context, _ string) ([]*order.Order, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.orders, nil
}
func (m *mockOrderRepo) UpdateStatus(_ context.Context, id polyid.OrderID, _ order.OrderStatus) error {
	m.updated = append(m.updated, id)
	return nil
}

// --- helpers ---

func openState() windowstate.WindowState {
	return windowstate.WindowState{
		Asset:    "btc",
		MarketID: "market-1",
		Status:   windowstate.WindowOpen,
	}
}

func makeOrder(t *testing.T) *order.Order {
	t.Helper()
	o, err := order.New(order.Params{
		MarketID:   "market-1",
		TokenID:    polyid.TokenID("token-up"),
		Side:       order.Buy,
		Outcome:    market.Up,
		Price:      decimal.RequireFromString("0.60"),
		Size:       decimal.RequireFromString("10"),
		Type:       order.FOK,
		FeeRateBps: 100,
	})
	require.NoError(t, err)
	return o
}

// --- tests ---

func TestCloseWindow_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     dto.Input
		state     windowstate.WindowState
		getErr    error
		listErr   error
		saveErr   error
		orders    []*order.Order
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output, store *mockWindowStore, repo *mockOrderRepo)
	}{
		{
			name:      "empty asset returns client error",
			input:     dto.Input{Asset: ""},
			state:     openState(),
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "GetWindowState failure returns internal error",
			input:     dto.Input{Asset: "btc"},
			getErr:    errors.New("store down"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
		{
			name:      "window not open returns client error",
			input:     dto.Input{Asset: "btc"},
			state:     windowstate.WindowState{Asset: "btc", Status: windowstate.WindowClosed},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:   "no open orders closes window successfully",
			input:  dto.Input{Asset: "btc"},
			state:  openState(),
			orders: []*order.Order{},
			checkOut: func(t *testing.T, out dto.Output, store *mockWindowStore, _ *mockOrderRepo) {
				t.Helper()
				assert.Equal(t, "btc", out.Asset)
				assert.Equal(t, 0, out.OrdersExpired)
				require.NotNil(t, store.saved)
				assert.Equal(t, windowstate.WindowClosed, store.saved.Status)
			},
		},
		{
			name:  "two open orders expire both and close window",
			input: dto.Input{Asset: "btc"},
			state: openState(),
			orders: func() []*order.Order {
				o1 := makeOrder(t)
				o2 := makeOrder(t)
				return []*order.Order{o1, o2}
			}(),
			checkOut: func(t *testing.T, out dto.Output, store *mockWindowStore, repo *mockOrderRepo) {
				t.Helper()
				assert.Equal(t, "btc", out.Asset)
				assert.Equal(t, 2, out.OrdersExpired)
				require.NotNil(t, store.saved)
				assert.Equal(t, windowstate.WindowClosed, store.saved.Status)
				assert.Len(t, repo.updated, 2)
			},
		},
		{
			name:      "SaveWindowState failure returns internal error",
			input:     dto.Input{Asset: "btc"},
			state:     openState(),
			orders:    []*order.Order{},
			saveErr:   errors.New("save failed"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := &mockWindowStore{state: tt.state, getErr: tt.getErr, saveErr: tt.saveErr}
			repo := &mockOrderRepo{orders: tt.orders, listErr: tt.listErr}
			uc := closewindow.New(store, repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				if tt.checkOut != nil {
					tt.checkOut(t, out, store, repo)
				}
			}
		})
	}
}
