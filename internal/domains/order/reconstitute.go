package order

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// ReconstitutedParams holds all fields needed to reconstitute an Order from storage.
type ReconstitutedParams struct {
	ID            polyid.OrderID
	MarketID      string
	TokenID       polyid.TokenID
	Side          Side
	Outcome       market.Outcome
	Price         decimal.Decimal
	Size          decimal.Decimal
	Type          OrderType
	Expiration    time.Time
	FeeRateBps    uint64
	SignatureType uint8
	Status        OrderStatus
	CreatedAt     time.Time
}

// Reconstitute rebuilds an Order from persisted state without running factory validation.
// Use only in repository scan functions — never in application code.
func Reconstitute(p ReconstitutedParams) *Order {
	return &Order{
		id:            p.ID,
		marketID:      p.MarketID,
		tokenID:       p.TokenID,
		side:          p.Side,
		outcome:       p.Outcome,
		price:         p.Price,
		size:          p.Size,
		orderType:     p.Type,
		expiration:    p.Expiration,
		feeRateBps:    p.FeeRateBps,
		signatureType: p.SignatureType,
		status:        p.Status,
		createdAt:     p.CreatedAt,
	}
}
