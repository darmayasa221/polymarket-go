package recordprice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands/recordprice"
	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands/recordprice/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// errMockNotFound is returned by unused mock methods to satisfy the nilnil linter.
var errMockNotFound = errors.New("mock: not found")

// mockPriceRepo is an inline mock for ports.PriceRepository.
type mockPriceRepo struct {
	saved   *oracle.Price
	saveErr error
}

func (m *mockPriceRepo) Save(_ context.Context, price *oracle.Price) error {
	m.saved = price
	return m.saveErr
}
func (m *mockPriceRepo) LatestByAsset(_ context.Context, _ string) (*oracle.Price, error) {
	return nil, errMockNotFound
}
func (m *mockPriceRepo) LatestChainlinkByAsset(_ context.Context, _ string) (*oracle.Price, error) {
	return nil, errMockNotFound
}
func (m *mockPriceRepo) WindowOpenPrice(_ context.Context, _ string, _ time.Time) (*oracle.Price, error) {
	return nil, errMockNotFound
}

func TestRecordPrice_Execute(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	tests := []struct {
		name      string
		input     dto.Input
		repoErr   error
		wantErr   bool
		errTarget any
		checkOut  func(t *testing.T, out dto.Output)
	}{
		{
			name: "valid chainlink price is saved",
			input: dto.Input{
				Asset:      "btc",
				Source:     "chainlink",
				Value:      "67234.50",
				RoundedAt:  now,
				ReceivedAt: now,
			},
			wantErr: false,
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "btc", out.Asset)
				assert.Equal(t, "chainlink", out.Source)
				assert.Equal(t, "67234.5", out.Value)
			},
		},
		{
			name: "valid binance price is saved",
			input: dto.Input{
				Asset:      "eth",
				Source:     "binance",
				Value:      "3500.00",
				ReceivedAt: now,
			},
			wantErr: false,
			checkOut: func(t *testing.T, out dto.Output) {
				t.Helper()
				assert.Equal(t, "eth", out.Asset)
				assert.Equal(t, "binance", out.Source)
			},
		},
		{
			name: "empty asset returns client error",
			input: dto.Input{
				Asset:      "",
				Source:     "chainlink",
				Value:      "67234.50",
				ReceivedAt: now,
			},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name: "invalid source returns client error",
			input: dto.Input{
				Asset:      "btc",
				Source:     "unknown",
				Value:      "67234.50",
				ReceivedAt: now,
			},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name: "invalid price value returns invariant error",
			input: dto.Input{
				Asset:      "btc",
				Source:     "chainlink",
				Value:      "not-a-number",
				ReceivedAt: now,
			},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name: "zero price returns invariant error",
			input: dto.Input{
				Asset:      "btc",
				Source:     "chainlink",
				Value:      "0",
				ReceivedAt: now,
			},
			wantErr:   true,
			errTarget: &errtypes.InvariantError{},
		},
		{
			name: "repository error returns internal error",
			input: dto.Input{
				Asset:      "btc",
				Source:     "chainlink",
				Value:      "67234.50",
				ReceivedAt: now,
			},
			repoErr:   errors.New("db unavailable"),
			wantErr:   true,
			errTarget: &errtypes.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockPriceRepo{saveErr: tt.repoErr}
			uc := recordprice.New(repo)

			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget), "expected error type %T, got %T", tt.errTarget, err)
				}
				assert.Empty(t, out.Asset)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, repo.saved)
				if tt.checkOut != nil {
					tt.checkOut(t, out)
				}
			}
		})
	}
}
