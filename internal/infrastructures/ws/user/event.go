// Package user provides the Polymarket user WebSocket event types and handler.
package user

// EventType classifies incoming user WebSocket messages.
type EventType string

const (
	// EventOrderFilled is emitted when an order is fully or partially matched.
	EventOrderFilled EventType = "order_filled"
	// EventOrderCanceled is emitted when an order is canceled.
	EventOrderCanceled EventType = "order_canceled"
)

// UserEvent is a typed event from the Polymarket user WebSocket.
type UserEvent struct {
	Type    EventType
	OrderID string
}
