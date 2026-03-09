package order

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Params holds all inputs required to construct an Order.
type Params struct {
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
}

// New creates and validates a new Order aggregate.
// This is the ONLY way to create an Order — never use struct literal.
func New(p Params) (*Order, error) {
	o := &Order{
		id:            polyid.OrderID(crypto.GenerateUUID()),
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
		status:        StatusOpen,
		createdAt:     timeutil.Now(),
	}
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return o, nil
}
