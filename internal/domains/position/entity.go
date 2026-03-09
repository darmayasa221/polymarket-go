// Package position defines the Position aggregate for tracking open outcome token holdings.
package position

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Position is the Position aggregate root.
// Represents outcome tokens currently held from one or more filled orders.
// Created only via New() — never use struct literal.
type Position struct {
	id       string
	asset    market.Asset
	tokenID  polyid.TokenID
	outcome  market.Outcome
	size     decimal.Decimal
	avgPrice decimal.Decimal
	marketID string
	openedAt time.Time
	closedAt *time.Time
}

// ID returns the position's local identifier.
func (p *Position) ID() string { return p.id }

// Asset returns the crypto asset.
func (p *Position) Asset() market.Asset { return p.asset }

// TokenID returns the ERC1155 outcome token ID.
func (p *Position) TokenID() polyid.TokenID { return p.tokenID }

// Outcome returns Up or Down.
func (p *Position) Outcome() market.Outcome { return p.outcome }

// Size returns the number of outcome shares held.
func (p *Position) Size() decimal.Decimal { return p.size }

// AvgPrice returns the weighted average entry price.
func (p *Position) AvgPrice() decimal.Decimal { return p.avgPrice }

// MarketID returns the associated market's Gamma event ID.
func (p *Position) MarketID() string { return p.marketID }

// OpenedAt returns when the position was first opened.
func (p *Position) OpenedAt() time.Time { return p.openedAt }

// ClosedAt returns when the position was closed, or nil if still open.
func (p *Position) ClosedAt() *time.Time { return p.closedAt }
