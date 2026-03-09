package order

import "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"

// Validate checks all business invariants for the Order.
func (o *Order) Validate() error {
	if o.marketID == "" {
		return types.NewInvariantError(ErrMarketIDRequired)
	}
	if o.tokenID.IsEmpty() {
		return types.NewInvariantError(ErrTokenIDRequired)
	}
	if o.price.IsZero() || o.price.IsNegative() {
		return types.NewInvariantError(ErrPriceInvalid)
	}
	if o.size.IsZero() || o.size.IsNegative() {
		return types.NewInvariantError(ErrSizeInvalid)
	}
	if o.orderType == GTD && o.expiration.IsZero() {
		return types.NewInvariantError(ErrExpirationRequired)
	}
	return nil
}
