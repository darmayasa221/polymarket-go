package ports

import (
	"context"
	"time"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

// PriceRepository stores and retrieves oracle.Price observations.
type PriceRepository interface {
	Save(ctx context.Context, price *oracle.Price) error
	LatestByAsset(ctx context.Context, asset string) (*oracle.Price, error)
	LatestChainlinkByAsset(ctx context.Context, asset string) (*oracle.Price, error)
	// WindowOpenPrice returns the first Chainlink price at or after windowStart.
	WindowOpenPrice(ctx context.Context, asset string, windowStart time.Time) (*oracle.Price, error)
}
