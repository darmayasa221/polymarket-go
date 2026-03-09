package oracle

import (
	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// PredictOutcome returns the predicted resolution based on open and current price.
// Returns Up if currentPrice >= openPrice, Down otherwise.
// This mirrors Polymarket's resolution rule: close >= open → "Up" wins.
func PredictOutcome(openPrice, currentPrice decimal.Decimal) market.Outcome {
	if currentPrice.GreaterThanOrEqual(openPrice) {
		return market.Up
	}
	return market.Down
}
