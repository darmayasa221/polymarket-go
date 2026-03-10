package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
)

const (
	userPingInterval = 10 * time.Second
	// UserEndpoint is the production user WebSocket URL.
	UserEndpoint = "wss://ws-subscriptions-clob.polymarket.com/ws/user"
)

// Handler connects to the Polymarket user WebSocket and emits UserEvent values.
type Handler struct{ url string }

// New creates a user Handler targeting the given WebSocket URL.
func New(url string) *Handler { return &Handler{url: url} }

// Start dials the user WebSocket, subscribes with L2 auth,
// and returns a read-only channel of UserEvent values.
func (h *Handler) Start(ctx context.Context, cfg clob.Config) (<-chan UserEvent, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, h.url, nil)
	if err != nil {
		return nil, fmt.Errorf("user ws: dial: %w", err)
	}
	// User channel auth uses raw credentials nested in an "auth" object.
	// See: https://docs.polymarket.com/developers/CLOB/websocket/user-channel
	sub := map[string]any{
		"type": "user",
		"auth": map[string]any{
			"apiKey":     cfg.APIKey,
			"secret":     cfg.APISecret,
			"passphrase": cfg.APIPassphrase,
		},
	}
	if err := conn.WriteJSON(sub); err != nil {
		conn.Close()
		return nil, fmt.Errorf("user ws: subscribe: %w", err)
	}
	out := make(chan UserEvent, 64)
	go readUserLoop(ctx, conn, out)
	return out, nil
}

// pumpMessages reads raw WebSocket frames and forwards them to msgCh until
// the connection closes or ctx is done.
func pumpMessages(ctx context.Context, conn *websocket.Conn, msgCh chan<- []byte) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			close(msgCh)
			return
		}
		select {
		case msgCh <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func readUserLoop(ctx context.Context, conn *websocket.Conn, out chan<- UserEvent) {
	defer close(out)
	defer conn.Close()

	pingTicker := time.NewTicker(userPingInterval)
	defer pingTicker.Stop()

	msgCh := make(chan []byte, 16)
	go pumpMessages(ctx, conn, msgCh)

	for {
		select {
		case <-ctx.Done():
			// defer conn.Close() will unblock pumpMessages' ReadMessage call.
			return
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			handleUserMessage(msg, out)
		}
	}
}

func handleUserMessage(raw []byte, out chan<- UserEvent) {
	// All user WS messages have top-level event_type, id, market, type, status.
	// See: https://docs.polymarket.com/market-data/websocket/user-channel
	var msg struct {
		EventType string `json:"event_type"`
		Type      string `json:"type"`   // "PLACEMENT", "CANCELLATION", "MATCHED", etc.
		Status    string `json:"status"` // trade status: "MATCHED", "CONFIRMED", etc.
		ID        string `json:"id"`
		Market    string `json:"market"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}
	ev := UserEvent{
		EventType: EventType(msg.EventType),
		OrderType: OrderType(msg.Type),
		Status:    TradeStatus(msg.Status),
		OrderID:   msg.ID,
		Market:    msg.Market,
	}
	switch ev.EventType {
	case EventTrade, EventOrder:
		select {
		case out <- ev:
		default:
		}
	}
}
