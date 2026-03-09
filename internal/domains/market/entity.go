package market

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
)

// Market is the Market aggregate root for a single 5-minute Up/Down crypto market.
// Created only via New() — never use struct literal.
type Market struct {
	id          string
	slug        polyid.SlugID
	asset       Asset
	windowStart time.Time
	conditionID polyid.ConditionID
	upTokenID   polyid.TokenID
	downTokenID polyid.TokenID
	tickSize    decimal.Decimal
	feeEnabled  bool
	active      bool
}

// ID returns the Gamma event ID.
func (m *Market) ID() string { return m.id }

// Slug returns the predictable market slug.
func (m *Market) Slug() polyid.SlugID { return m.slug }

// Asset returns the crypto asset (BTC/ETH/SOL/XRP).
func (m *Market) Asset() Asset { return m.asset }

// WindowStart returns the 5-minute window start time.
func (m *Market) WindowStart() time.Time { return m.windowStart }

// WindowEnd returns WindowStart + 5 minutes.
func (m *Market) WindowEnd() time.Time { return m.windowStart.Add(5 * time.Minute) }

// ConditionID returns the CTF condition identifier.
func (m *Market) ConditionID() polyid.ConditionID { return m.conditionID }

// UpTokenID returns the ERC1155 token ID for the "Up" outcome.
func (m *Market) UpTokenID() polyid.TokenID { return m.upTokenID }

// DownTokenID returns the ERC1155 token ID for the "Down" outcome.
func (m *Market) DownTokenID() polyid.TokenID { return m.downTokenID }

// TickSize returns the minimum price increment for this market.
func (m *Market) TickSize() decimal.Decimal { return m.tickSize }

// FeeEnabled returns whether fees are charged on orders.
func (m *Market) FeeEnabled() bool { return m.feeEnabled }

// Active returns whether the market accepts new orders.
func (m *Market) Active() bool { return m.active }

// TokenIDFor returns the TokenID for the given outcome.
func (m *Market) TokenIDFor(o Outcome) polyid.TokenID {
	if o == Up {
		return m.upTokenID
	}
	return m.downTokenID
}
