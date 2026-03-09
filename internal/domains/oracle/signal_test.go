package oracle_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

func TestPredictOutcome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		openPrice    decimal.Decimal
		currentPrice decimal.Decimal
		want         market.Outcome
	}{
		{
			name:         "close > open returns Up",
			openPrice:    decimal.NewFromFloat(60000),
			currentPrice: decimal.NewFromFloat(61000),
			want:         market.Up,
		},
		{
			name:         "close == open returns Up",
			openPrice:    decimal.NewFromFloat(60000),
			currentPrice: decimal.NewFromFloat(60000),
			want:         market.Up,
		},
		{
			name:         "close < open returns Down",
			openPrice:    decimal.NewFromFloat(60000),
			currentPrice: decimal.NewFromFloat(59000),
			want:         market.Down,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := oracle.PredictOutcome(tt.openPrice, tt.currentPrice)
			assert.Equal(t, tt.want, got)
		})
	}
}
