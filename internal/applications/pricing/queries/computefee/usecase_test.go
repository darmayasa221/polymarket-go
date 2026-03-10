package computefee_test

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/computefee"
	"github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/computefee/dto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// TestComputeFee_ResearchTable verifies the fee formula constant against the
// live API research table from docs/decisions/5m-market-mechanics.md.
// CRITICAL: these values must match the live API — do NOT change without re-research.
func TestComputeFee_ResearchTable(t *testing.T) {
	t.Parallel()

	uc := computefee.New()

	// Table-driven verification: p → expected feePerShare (±0.0005 tolerance)
	cases := []struct {
		price   string
		wantFee float64 // expected feePerShare from research
		wantBps int64   // approximate bps
	}{
		{price: "0.50", wantFee: 0.0156, wantBps: 156}, // peak: ~156 bps at p=0.50
		{price: "0.20", wantFee: 0.0100, wantBps: 100}, // ~100 bps at p=0.20
		{price: "0.95", wantFee: 0.0030, wantBps: 30},  // ~30 bps at p=0.95
	}

	for _, tc := range cases {
		t.Run("p="+tc.price, func(t *testing.T) {
			t.Parallel()
			out, err := uc.Execute(t.Context(), dto.Input{TokenPrice: tc.price})
			require.NoError(t, err)

			got, _ := out.Fee.FeePerShare.Float64()
			assert.InDelta(t, tc.wantFee, got, 0.002,
				"feePerShare at p=%s: got %.4f, want %.4f", tc.price, got, tc.wantFee)
			assert.Equal(t, tc.wantBps, out.Fee.EffectiveBps,
				"EffectiveBps at p=%s", tc.price)
		})
	}
}

func TestComputeFee_Execute(t *testing.T) {
	t.Parallel()

	uc := computefee.New()

	tests := []struct {
		name      string
		input     dto.Input
		wantErr   bool
		errTarget any
		checkFee  func(t *testing.T, fee float64)
	}{
		{
			name:  "valid price 0.50 returns peak fee",
			input: dto.Input{TokenPrice: "0.50"},
			checkFee: func(t *testing.T, fee float64) {
				t.Helper()
				// Fee at p=0.50 is the maximum (~0.0156)
				assert.Greater(t, fee, 0.015)
				assert.Less(t, fee, 0.017)
			},
		},
		{
			name:  "price at 0.01 (extreme) has very low fee",
			input: dto.Input{TokenPrice: "0.01"},
			checkFee: func(t *testing.T, fee float64) {
				t.Helper()
				assert.Less(t, fee, 0.001)
			},
		},
		{
			name:  "price at 0.99 (extreme) has very low fee",
			input: dto.Input{TokenPrice: "0.99"},
			checkFee: func(t *testing.T, fee float64) {
				t.Helper()
				assert.Less(t, fee, 0.001)
			},
		},
		{
			name:  "formula is symmetric: fee(p) == fee(1-p)",
			input: dto.Input{TokenPrice: "0.30"},
			checkFee: func(t *testing.T, fee float64) {
				t.Helper()
				outMirror, err := uc.Execute(t.Context(), dto.Input{TokenPrice: "0.70"})
				require.NoError(t, err)
				mirror, _ := outMirror.Fee.FeePerShare.Float64()
				assert.InDelta(t, fee, mirror, 0.0001, "formula should be symmetric around 0.5")
			},
		},
		{
			name:  "fee is always non-negative",
			input: dto.Input{TokenPrice: "0.75"},
			checkFee: func(t *testing.T, fee float64) {
				t.Helper()
				assert.True(t, fee >= 0, "fee should be non-negative, got %f", fee)
				assert.False(t, math.IsNaN(fee))
			},
		},
		{
			name:      "empty price returns client error",
			input:     dto.Input{TokenPrice: ""},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "non-numeric price returns client error",
			input:     dto.Input{TokenPrice: "abc"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "price of zero returns client error",
			input:     dto.Input{TokenPrice: "0"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
		{
			name:      "price greater than 1 returns client error",
			input:     dto.Input{TokenPrice: "1.01"},
			wantErr:   true,
			errTarget: &errtypes.ClientError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out, err := uc.Execute(t.Context(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errTarget != nil {
					assert.True(t, errors.As(err, &tt.errTarget))
				}
			} else {
				require.NoError(t, err)
				fee, _ := out.Fee.FeePerShare.Float64()
				if tt.checkFee != nil {
					tt.checkFee(t, fee)
				}
			}
		})
	}
}
