package windowstate

import (
	"time"

	"github.com/shopspring/decimal"
)

// WindowStatus represents the lifecycle state of a 5-minute trading window.
type WindowStatus string

const (
	// WindowOpen means the window is active and orders can be placed.
	WindowOpen WindowStatus = "open"
	// WindowClosed means the window has ended and orders are no longer accepted.
	WindowClosed WindowStatus = "closed"
	// WindowSettling means the window is closed and awaiting Chainlink resolution (~2-3 min).
	WindowSettling WindowStatus = "settling"
)

// OrderSummary is a read-only snapshot of an order embedded in WindowState.
type OrderSummary struct {
	OrderID string
	Side    string // "buy" | "sell"
	Outcome string // "Up" | "Down"
	Price   string
	Size    string
	Status  string
}

// WindowState is the DTO returned by GetWindowState query.
// It is ephemeral — not persisted between bot restarts.
type WindowState struct {
	MarketID    string
	Asset       string
	WindowStart time.Time
	WindowEnd   time.Time
	ConditionID string
	UpTokenID   string
	DownTokenID string
	TickSize    string
	OpenPrice   decimal.Decimal
	Status      WindowStatus
	OpenOrders  []OrderSummary
}
