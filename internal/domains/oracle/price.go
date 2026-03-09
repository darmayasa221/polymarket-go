package oracle

import (
	"time"

	"github.com/shopspring/decimal"
)

// Price represents a single price observation from an external feed.
// Created only via New() — never use struct literal.
type Price struct {
	asset      string
	source     PriceSource
	value      decimal.Decimal
	roundedAt  time.Time // Chainlink round timestamp (zero for Binance)
	receivedAt time.Time
}

// Asset returns the crypto ticker (btc/eth/sol/xrp).
func (p *Price) Asset() string { return p.asset }

// Source returns where this price came from.
func (p *Price) Source() PriceSource { return p.source }

// Value returns the price in USD.
func (p *Price) Value() decimal.Decimal { return p.value }

// RoundedAt returns the Chainlink round timestamp (zero for Binance prices).
func (p *Price) RoundedAt() time.Time { return p.roundedAt }

// ReceivedAt returns when this price was received by the bot.
func (p *Price) ReceivedAt() time.Time { return p.receivedAt }
