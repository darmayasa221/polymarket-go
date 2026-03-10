package position

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// ReconstitutedParams holds all fields needed to reconstitute a Position from storage.
type ReconstitutedParams struct {
	ID       string
	Asset    market.Asset
	TokenID  polyid.TokenID
	Outcome  market.Outcome
	Size     decimal.Decimal
	AvgPrice decimal.Decimal
	MarketID string
	OpenedAt time.Time
	ClosedAt *time.Time
}

// Reconstitute rebuilds a Position from persisted state without running factory validation.
// Use only in repository scan functions — never in application code.
func Reconstitute(p ReconstitutedParams) *Position {
	return &Position{
		id:       p.ID,
		asset:    p.Asset,
		tokenID:  p.TokenID,
		outcome:  p.Outcome,
		size:     p.Size,
		avgPrice: p.AvgPrice,
		marketID: p.MarketID,
		openedAt: p.OpenedAt,
		closedAt: p.ClosedAt,
	}
}
