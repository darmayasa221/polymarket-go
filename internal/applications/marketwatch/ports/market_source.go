package ports

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// MarketSource fetches active 5-minute markets from the Polymarket Gamma API.
type MarketSource interface {
	FetchActive5mMarkets(ctx context.Context) ([]*market.Market, error)
}
