// Package user provides the Polymarket user WebSocket event types and handler.
package user

// EventType classifies incoming user WebSocket messages by event_type field.
type EventType string

const (
	// EventTrade is emitted when an order is matched or a trade status changes.
	EventTrade EventType = "trade"
	// EventOrder is emitted when an order is placed, updated, or canceled.
	EventOrder EventType = "order"
)

// OrderType is the type field inside an order event.
type OrderType string

const (
	OrderPlacement    OrderType = "PLACEMENT"
	OrderCancellation OrderType = "CANCELLATION"
)

// TradeStatus is the status field inside a trade event.
type TradeStatus string

const (
	TradeMatched   TradeStatus = "MATCHED"
	TradeConfirmed TradeStatus = "CONFIRMED"
)

// UserEvent is a typed event from the Polymarket user WebSocket.
type UserEvent struct {
	EventType EventType
	OrderType OrderType
	Status    TradeStatus
	OrderID   string
	Market    string
}
