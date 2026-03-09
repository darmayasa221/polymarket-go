package market

import (
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// Validate checks all business invariants for the Market.
func (m *Market) Validate() error {
	if m.id == "" {
		return types.NewInvariantError(ErrIDRequired)
	}
	if !m.asset.IsValid() {
		return types.NewInvariantError(ErrInvalidAsset)
	}
	if m.windowStart.IsZero() {
		return types.NewInvariantError(ErrWindowStartRequired)
	}
	if m.conditionID.IsEmpty() {
		return types.NewInvariantError(ErrConditionIDRequired)
	}
	if m.upTokenID.IsEmpty() {
		return types.NewInvariantError(ErrUpTokenRequired)
	}
	if m.downTokenID.IsEmpty() {
		return types.NewInvariantError(ErrDownTokenRequired)
	}
	if m.tickSize.IsZero() || m.tickSize.IsNegative() {
		return types.NewInvariantError(ErrTickSizeInvalid)
	}
	return nil
}
