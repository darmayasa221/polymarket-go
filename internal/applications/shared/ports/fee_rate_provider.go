package ports

import "context"

// FeeRateProvider fetches the live fee rate from the CLOB.
// Defined here (shared) so both pricing and trading contexts can use it without duplication.
type FeeRateProvider interface {
	FetchFeeRate(ctx context.Context, tokenID string) (uint64, error)
}
