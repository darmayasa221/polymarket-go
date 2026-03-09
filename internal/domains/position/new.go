package position

import (
	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Params holds all inputs required to construct a Position.
type Params struct {
	Asset    market.Asset
	TokenID  polyid.TokenID
	Outcome  market.Outcome
	Size     decimal.Decimal
	AvgPrice decimal.Decimal
	MarketID string
}

// New creates and validates a new Position aggregate.
// This is the ONLY way to create a Position — never use struct literal.
func New(p Params) (*Position, error) {
	pos := &Position{
		id:       crypto.GenerateUUID(),
		asset:    p.Asset,
		tokenID:  p.TokenID,
		outcome:  p.Outcome,
		size:     p.Size,
		avgPrice: p.AvgPrice,
		marketID: p.MarketID,
		openedAt: timeutil.Now(),
		closedAt: nil,
	}
	if err := pos.Validate(); err != nil {
		return nil, err
	}
	return pos, nil
}
