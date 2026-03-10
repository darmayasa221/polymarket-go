// Package rtds implements the Polymarket RTDS WebSocket handler for live price feeds.
package rtds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

const (
	rtdsPingInterval = 5 * time.Second
	// RTDSEndpoint is the production RTDS WebSocket URL.
	RTDSEndpoint = "wss://ws-live-data.polymarket.com"
)

// rtdsSubscription is the message sent immediately after connecting.
// Subscribes to Chainlink and Binance price feeds for all tracked assets.
var rtdsSubscription = map[string]any{
	"action": "subscribe",
	"subscriptions": []map[string]any{
		{"topic": "crypto_prices_chainlink", "type": "*"},
		{"topic": "crypto_prices", "type": "*"},
	},
}

// Handler connects to the RTDS WebSocket and emits oracle.Price observations.
type Handler struct{ url string }

// New creates an RTDS Handler targeting the given WebSocket URL.
// Use RTDSEndpoint for production; use a mock server URL in tests.
func New(url string) *Handler { return &Handler{url: url} }

// Start dials the RTDS WebSocket, subscribes to price feeds,
// and returns a read-only channel of oracle.Price values.
// The channel is closed when ctx is canceled or the connection is lost.
func (h *Handler) Start(ctx context.Context) (<-chan *oracle.Price, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, h.url, nil)
	if err != nil {
		return nil, fmt.Errorf("rtds: dial: %w", err)
	}

	// Send subscription message immediately after connecting.
	if err := conn.WriteJSON(rtdsSubscription); err != nil {
		conn.Close()
		return nil, fmt.Errorf("rtds: subscribe: %w", err)
	}

	out := make(chan *oracle.Price, 64)
	go h.readLoop(ctx, conn, out)
	return out, nil
}

// rtdsEnvelope is the top-level RTDS message structure.
type rtdsEnvelope struct {
	Topic   string          `json:"topic"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// rtdsPayload is the inner payload for both Chainlink and Binance messages.
type rtdsPayload struct {
	Symbol    string  `json:"symbol"`
	Timestamp int64   `json:"timestamp"` // Unix milliseconds
	Value     float64 `json:"value"`
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

func (h *Handler) readLoop(ctx context.Context, conn *websocket.Conn, out chan<- *oracle.Price) {
	defer close(out)
	defer conn.Close()

	pingTicker := time.NewTicker(rtdsPingInterval)
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
			h.handleMessage(msg, out)
		}
	}
}

func (h *Handler) handleMessage(raw []byte, out chan<- *oracle.Price) {
	var env rtdsEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return
	}

	var source oracle.PriceSource
	switch env.Topic {
	case "crypto_prices_chainlink":
		source = oracle.SourceChainlink
	case "crypto_prices":
		source = oracle.SourceBinance
	default:
		return // ignore subscription confirmations and other control messages
	}

	var pd rtdsPayload
	if err := json.Unmarshal(env.Payload, &pd); err != nil {
		return
	}

	asset := symbolToAsset(pd.Symbol, env.Topic)
	if asset == "" {
		return
	}

	params := oracle.Params{
		Asset:      asset,
		Source:     source,
		Value:      decimal.NewFromFloat(pd.Value),
		ReceivedAt: timeutil.Now(),
	}
	// RoundedAt is the Chainlink round timestamp (Unix milliseconds).
	if pd.Timestamp > 0 {
		params.RoundedAt = time.UnixMilli(pd.Timestamp).UTC()
	}

	price, err := oracle.New(params)
	if err != nil {
		log.Printf("rtds: invalid price observation: %v", err)
		return
	}

	select {
	case out <- price:
	default: // drop if consumer is slow — never block
	}
}

// symbolToAsset maps a raw RTDS symbol to a normalised asset name.
// Chainlink uses "btc/usd" format; Binance uses "btcusdt" format.
func symbolToAsset(symbol, topic string) string {
	symbol = strings.ToLower(symbol)
	if topic == "crypto_prices_chainlink" {
		// "btc/usd" → "btc"
		parts := strings.SplitN(symbol, "/", 2)
		return parts[0]
	}
	// Binance: "btcusdt" → "btc", "ethusdt" → "eth", etc.
	for _, known := range []string{"btc", "eth", "sol", "xrp"} {
		if strings.HasPrefix(symbol, known) {
			return known
		}
	}
	return ""
}
