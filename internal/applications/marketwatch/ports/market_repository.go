package ports

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// MarketRepository persists and retrieves Market aggregates.
type MarketRepository interface {
	Save(ctx context.Context, m *market.Market) error
	FindByAssetAndWindow(ctx context.Context, asset market.Asset, windowStart time.Time) (*market.Market, error)
	UpdateTickSize(ctx context.Context, conditionID string, newTickSize decimal.Decimal) error
	ListActive(ctx context.Context) ([]*market.Market, error)
}
