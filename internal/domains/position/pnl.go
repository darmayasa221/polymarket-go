package position

import "github.com/shopspring/decimal"

// UnrealisedPnL returns the profit/loss if the position were closed at currentPrice.
// Formula: size * (currentPrice - avgPrice).
func (p *Position) UnrealisedPnL(currentPrice decimal.Decimal) decimal.Decimal {
	return p.size.Mul(currentPrice.Sub(p.avgPrice))
}

// RealisedPnL returns the profit/loss from closing the position at exitPrice.
// Formula: size * (exitPrice - avgPrice).
func (p *Position) RealisedPnL(exitPrice decimal.Decimal) decimal.Decimal {
	return p.size.Mul(exitPrice.Sub(p.avgPrice))
}
