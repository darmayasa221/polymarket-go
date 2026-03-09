package oracle

import "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"

// Validate checks all invariants for a Price observation.
func (p *Price) Validate() error {
	if p.asset == "" {
		return types.NewInvariantError(ErrAssetRequired)
	}
	if !p.source.IsValid() {
		return types.NewInvariantError(ErrInvalidSource)
	}
	if p.value.IsZero() || p.value.IsNegative() {
		return types.NewInvariantError(ErrPriceValueInvalid)
	}
	if p.receivedAt.IsZero() {
		return types.NewInvariantError(ErrReceivedAtRequired)
	}
	return nil
}
