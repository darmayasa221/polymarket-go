package market

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/slug"
)

// Params holds all inputs required to construct a Market.
type Params struct {
	ID          string
	Asset       Asset
	WindowStart time.Time
	ConditionID polyid.ConditionID
	UpTokenID   polyid.TokenID
	DownTokenID polyid.TokenID
	TickSize    decimal.Decimal
	FeeEnabled  bool
}

// New creates and validates a new Market aggregate.
// This is the ONLY way to create a Market — never use struct literal.
func New(p Params) (*Market, error) {
	m := &Market{
		id:          p.ID,
		slug:        slug.ForAsset(string(p.Asset), p.WindowStart),
		asset:       p.Asset,
		windowStart: p.WindowStart.UTC(),
		conditionID: p.ConditionID,
		upTokenID:   p.UpTokenID,
		downTokenID: p.DownTokenID,
		tickSize:    p.TickSize,
		feeEnabled:  p.FeeEnabled,
		active:      true,
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}
