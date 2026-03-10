package market

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/slug"
)

// ReconstitutedParams holds all fields needed to reconstitute a Market from storage.
type ReconstitutedParams struct {
	ID          string
	Asset       Asset
	WindowStart time.Time
	ConditionID polyid.ConditionID
	UpTokenID   polyid.TokenID
	DownTokenID polyid.TokenID
	TickSize    decimal.Decimal
	FeeEnabled  bool
	Active      bool
}

// Reconstitute rebuilds a Market from persisted state without running factory validation.
// Use only in repository scan functions — never in application code.
func Reconstitute(p ReconstitutedParams) *Market {
	return &Market{
		id:          p.ID,
		slug:        slug.ForAsset(string(p.Asset), p.WindowStart),
		asset:       p.Asset,
		windowStart: p.WindowStart.UTC(),
		conditionID: p.ConditionID,
		upTokenID:   p.UpTokenID,
		downTokenID: p.DownTokenID,
		tickSize:    p.TickSize,
		feeEnabled:  p.FeeEnabled,
		active:      p.Active,
	}
}
