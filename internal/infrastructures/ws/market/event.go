// Package market provides the Polymarket market WebSocket event types and handler.
package market

// EventType classifies incoming market WebSocket messages.
type EventType string

const (
	// EventTickSizeChange is emitted when the minimum tick size changes mid-window.
	EventTickSizeChange EventType = "tick_size_change"
	// EventNewMarket is emitted when a new 5-minute window opens.
	EventNewMarket EventType = "new_market"
	// EventMarketResolved is emitted when a window settles.
	EventMarketResolved EventType = "market_resolved"
)

// TickSizeChangePayload carries the data for a tick_size_change event.
type TickSizeChangePayload struct {
	// ConditionID identifies the market (maps to condition_id in the DB).
	ConditionID string
	// NewTickSize is the updated tick size as a decimal string (e.g. "0.001").
	NewTickSize string
}

// MarketEvent is a typed event from the Polymarket market WebSocket.
type MarketEvent struct {
	Type           EventType
	TickSizeChange *TickSizeChangePayload // non-nil only when Type == EventTickSizeChange
}
