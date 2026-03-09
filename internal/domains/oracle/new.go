package oracle

import (
	"time"

	"github.com/shopspring/decimal"
)

// Params holds inputs to create a new Price observation.
type Params struct {
	Asset      string
	Source     PriceSource
	Value      decimal.Decimal
	RoundedAt  time.Time
	ReceivedAt time.Time
}

// New creates and validates a new Price observation.
func New(p Params) (*Price, error) {
	price := &Price{
		asset:      p.Asset,
		source:     p.Source,
		value:      p.Value,
		roundedAt:  p.RoundedAt,
		receivedAt: p.ReceivedAt,
	}
	if err := price.Validate(); err != nil {
		return nil, err
	}
	return price, nil
}
