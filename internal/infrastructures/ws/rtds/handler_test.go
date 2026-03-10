package rtds_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/rtds"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func TestRTDSHandler_Start_ReceivesChainlinkPrice(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// RTDS wire format: top-level "topic" and "payload" fields.
		// Chainlink symbol format: "btc/usd" (slash-separated pair).
		inner, _ := json.Marshal(map[string]any{
			"symbol":    "btc/usd",
			"timestamp": time.Now().UnixMilli(),
			"value":     65000.50,
		})
		msg, _ := json.Marshal(map[string]any{
			"topic":   "crypto_prices_chainlink",
			"type":    "*",
			"payload": json.RawMessage(inner),
		})
		_ = conn.WriteMessage(websocket.TextMessage, msg)

		// Keep alive until client disconnects.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	handler := rtds.New(wsURL)

	ctx := t.Context()
	ch, err := handler.Start(ctx)
	require.NoError(t, err)

	select {
	case price := <-ch:
		require.NotNil(t, price)
		assert.Equal(t, "btc", price.Asset())
		assert.Equal(t, "65000.5", price.Value().String())
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for price")
	}
}
