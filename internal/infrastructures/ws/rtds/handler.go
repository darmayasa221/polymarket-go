// Package rtds implements the Polymarket RTDS WebSocket handler for live price feeds.
package rtds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

const (
	rtdsPingInterval = 5 * time.Second
	// RTDSEndpoint is the production RTDS WebSocket URL.
	RTDSEndpoint = "wss://ws-live-data.polymarket.com"
)

// Handler connects to the RTDS WebSocket and emits oracle.Price observations.
type Handler struct{ url string }

// New creates an RTDS Handler targeting the given WebSocket URL.
// Use RTDSEndpoint for production; use a mock server URL in tests.
func New(url string) *Handler { return &Handler{url: url} }

// Start dials the RTDS WebSocket and returns a read-only channel of oracle.Price values.
// The channel is closed when ctx is canceled or the connection is lost.
func (h *Handler) Start(ctx context.Context) (<-chan *oracle.Price, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, h.url, nil)
	if err != nil {
		return nil, fmt.Errorf("rtds: dial: %w", err)
	}
	out := make(chan *oracle.Price, 64)
	go h.readLoop(ctx, conn, out)
	return out, nil
}

type rtdsEnvelope struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

type rtdsPriceData struct {
	Asset     string `json:"asset"`
	Price     string `json:"price"`
	RoundedAt string `json:"rounded_at"`
}

func (h *Handler) readLoop(ctx context.Context, conn *websocket.Conn, out chan<- *oracle.Price) {
	defer close(out)
	defer conn.Close()

	pingTicker := time.NewTicker(rtdsPingInterval)
	defer pingTicker.Stop()

	msgCh := make(chan []byte, 16)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				close(msgCh)
				return
			}
			msgCh <- msg
		}
	}()

	for {
		select {
		case <-ctx.Done():
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
	switch env.EventType {
	case "crypto_prices_chainlink":
		source = oracle.SourceChainlink
	case "crypto_prices_binance":
		source = oracle.SourceBinance
	default:
		return
	}
	var pd rtdsPriceData
	if err := json.Unmarshal(env.Data, &pd); err != nil {
		return
	}
	value, err := decimal.NewFromString(pd.Price)
	if err != nil {
		return
	}
	params := oracle.Params{
		Asset:      pd.Asset,
		Source:     source,
		Value:      value,
		ReceivedAt: time.Now().UTC(),
	}
	if pd.RoundedAt != "" {
		if t, parseErr := time.Parse(time.RFC3339, pd.RoundedAt); parseErr == nil {
			params.RoundedAt = t
		}
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
