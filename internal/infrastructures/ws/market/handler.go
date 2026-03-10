package market

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	marketPingInterval = 10 * time.Second
	// MarketEndpoint is the production market WebSocket URL.
	MarketEndpoint = "wss://ws-subscriptions-clob.polymarket.com/ws/market"
)

// Handler connects to the Polymarket market WebSocket and emits MarketEvent values.
type Handler struct{ url string }

// New creates a market Handler targeting the given WebSocket URL.
func New(url string) *Handler { return &Handler{url: url} }

// Start dials the market WebSocket, subscribes with custom_feature_enabled=true,
// and returns a read-only channel of MarketEvent values.
// assetIDs is the list of condition IDs to subscribe to.
func (h *Handler) Start(ctx context.Context, assetIDs []string) (<-chan MarketEvent, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, h.url, nil)
	if err != nil {
		return nil, fmt.Errorf("market ws: dial: %w", err)
	}
	sub := map[string]any{
		"assets_ids":             assetIDs,
		"type":                   "market",
		"custom_feature_enabled": true,
	}
	if err := conn.WriteJSON(sub); err != nil {
		conn.Close()
		return nil, fmt.Errorf("market ws: subscribe: %w", err)
	}
	out := make(chan MarketEvent, 64)
	go readMarketLoop(ctx, conn, out)
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

func readMarketLoop(ctx context.Context, conn *websocket.Conn, out chan<- MarketEvent) {
	defer close(out)
	defer conn.Close()

	pingTicker := time.NewTicker(marketPingInterval)
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
			handleMarketMessage(msg, out)
		}
	}
}

func handleMarketMessage(raw []byte, out chan<- MarketEvent) {
	var env struct {
		EventType string          `json:"event_type"`
		Data      json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return
	}
	switch EventType(env.EventType) {
	case EventTickSizeChange:
		var payload struct {
			AssetID string `json:"asset_id"`
			NewSize string `json:"new_size"`
		}
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			return
		}
		select {
		case out <- MarketEvent{
			Type: EventTickSizeChange,
			TickSizeChange: &TickSizeChangePayload{
				ConditionID: payload.AssetID,
				NewTickSize: payload.NewSize,
			},
		}:
		default:
		}
	case EventNewMarket, EventMarketResolved:
		select {
		case out <- MarketEvent{Type: EventType(env.EventType)}:
		default:
		}
	}
}
