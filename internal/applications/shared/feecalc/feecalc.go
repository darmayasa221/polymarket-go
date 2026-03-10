package feecalc

import "github.com/shopspring/decimal"

// FeeResult is the output of ComputeFee query.
type FeeResult struct {
	TokenPrice   decimal.Decimal // input p (0.00–1.00)
	FeePerShare  decimal.Decimal // result of parabolic formula
	EffectiveBps int64           // approximate bps for logging/display
}
