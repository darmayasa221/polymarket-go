package order

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Order is the Order aggregate root for a single CLOB limit order.
// Created only via New() — never use struct literal.
type Order struct {
	id            polyid.OrderID
	marketID      string
	tokenID       polyid.TokenID
	side          Side
	outcome       market.Outcome
	price         decimal.Decimal
	size          decimal.Decimal
	orderType     OrderType
	expiration    time.Time
	feeRateBps    uint64
	signatureType uint8
	status        OrderStatus
	createdAt     time.Time
}

// ID returns the order's CLOB identifier.
func (o *Order) ID() polyid.OrderID { return o.id }

// MarketID returns the Gamma event ID of the market.
func (o *Order) MarketID() string { return o.marketID }

// TokenID returns the outcome token being traded.
func (o *Order) TokenID() polyid.TokenID { return o.tokenID }

// Side returns Buy or Sell.
func (o *Order) Side() Side { return o.side }

// Outcome returns Up or Down.
func (o *Order) Outcome() market.Outcome { return o.outcome }

// Price returns the limit price (0..1 range for binary markets).
func (o *Order) Price() decimal.Decimal { return o.price }

// Size returns the number of outcome shares.
func (o *Order) Size() decimal.Decimal { return o.size }

// Type returns the order's time-in-force type.
func (o *Order) Type() OrderType { return o.orderType }

// Expiration returns when a GTD order expires.
func (o *Order) Expiration() time.Time { return o.expiration }

// FeeRateBps returns the fee rate in basis points fetched from CLOB.
func (o *Order) FeeRateBps() uint64 { return o.feeRateBps }

// SignatureType returns 0 for EOA, 1 for Safe.
func (o *Order) SignatureType() uint8 { return o.signatureType }

// Status returns the current order status.
func (o *Order) Status() OrderStatus { return o.status }

// CreatedAt returns when the order was created locally.
func (o *Order) CreatedAt() time.Time { return o.createdAt }
