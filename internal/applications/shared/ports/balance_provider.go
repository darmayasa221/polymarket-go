package ports

import (
	"context"

	"github.com/shopspring/decimal"
)

// BalanceProvider fetches the current USDC.e collateral balance from the CLOB.
type BalanceProvider interface {
	FetchBalance(ctx context.Context) (decimal.Decimal, error)
}
