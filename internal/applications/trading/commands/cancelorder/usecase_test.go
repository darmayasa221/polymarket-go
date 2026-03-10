package cancelorder_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder"
	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// errMockNotFound satisfies nilnil linter.
var errMockNotFound = errors.New("mock: not found")

type mockOrderRepo struct {
	updateErr error
}

func (m *mockOrderRepo) Save(_ context.Context, _ *order.Order) error {
	return errors.New("mock: Save not used")
}
func (m *mockOrderRepo) FindByID(_ context.Context, _ polyid.OrderID) (*order.Order, error) {
	return nil, errMockNotFound
}
func (m *mockOrderRepo) ListOpenByMarket(_ context.Context, _ string) ([]*order.Order, error) {
	return nil, errMockNotFound
}
func (m *mockOrderRepo) UpdateStatus(_ context.Context, _ polyid.OrderID, _ order.OrderStatus) error {
	return m.updateErr
}

type mockOrderSubmitter struct {
	cancelErr error
	canceled  string
}

func (m *mockOrderSubmitter) Submit(_ context.Context, _ *order.Order, _ []byte) (string, error) {
	return "", errors.New("mock: Submit not used")
}
func (m *mockOrderSubmitter) Cancel(_ context.Context, clobOrderID string) error {
	m.canceled = clobOrderID
	return m.cancelErr
}

func TestCancelOrder_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     dto.Input
		cancelErr error
		updateErr error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output, sub *mockOrderSubmitter)
	}{
		{
			name:  "valid cancel succeeds",
			input: dto.Input{OrderID: "local-uuid", ClobOrderID: "clob-123"},
			checkOut: func(t *testing.T, out dto.Output, sub *mockOrderSubmitter) {
				t.Helper()
				assert.Equal(t, "local-uuid", out.OrderID)
				assert.Equal(t, "clob-123", sub.canceled)
			},
		},
		{
			name:      "empty order ID returns client error",
			input:     dto.Input{OrderID: "", ClobOrderID: "clob-123"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "empty CLOB order ID returns client error",
			input:     dto.Input{OrderID: "local-uuid", ClobOrderID: ""},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "CLOB cancel failure returns internal error",
			input:     dto.Input{OrderID: "local-uuid", ClobOrderID: "clob-123"},
			cancelErr: errors.New("CLOB rejected cancel"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockOrderRepo{updateErr: tt.updateErr}
			sub := &mockOrderSubmitter{cancelErr: tt.cancelErr}
			uc := cancelorder.New(repo, sub)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				if tt.checkOut != nil {
					tt.checkOut(t, out, sub)
				}
			}
		})
	}
}
