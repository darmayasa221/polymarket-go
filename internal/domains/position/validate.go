package position

import "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"

// Validate checks all business invariants for a Position.
func (p *Position) Validate() error {
	if p.marketID == "" {
		return types.NewInvariantError(ErrMarketIDRequired)
	}
	if p.tokenID.IsEmpty() {
		return types.NewInvariantError(ErrTokenIDRequired)
	}
	if p.size.IsZero() || p.size.IsNegative() {
		return types.NewInvariantError(ErrSizeInvalid)
	}
	if p.avgPrice.IsZero() || p.avgPrice.IsNegative() {
		return types.NewInvariantError(ErrAvgPriceInvalid)
	}
	return nil
}
