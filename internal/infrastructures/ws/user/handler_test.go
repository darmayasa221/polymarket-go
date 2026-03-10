package user_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
	wsuser "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/user"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func validCLOBConfig(baseURL string) clob.Config {
	return clob.Config{
		BaseURL:       baseURL,
		Address:       "0xAddress",
		APIKey:        "test-api-key",
		APISecret:     base64.StdEncoding.EncodeToString(make([]byte, 32)),
		APIPassphrase: "test-passphrase",
	}
}

func TestUserHandler_Start_ReceivesOrderCancellation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// read subscription message
		_, _, _ = conn.ReadMessage()

		// Send a real Polymarket user WS message:
		// id and market are top-level, event_type is "order", type is "CANCELLATION".
		msg, _ := json.Marshal(map[string]any{
			"event_type": "order",
			"type":       "CANCELLATION",
			"id":         "order-abc",
			"market":     "0xdeadbeef",
		})
		_ = conn.WriteMessage(websocket.TextMessage, msg)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	handler := wsuser.New(wsURL)
	cfg := validCLOBConfig(srv.URL)

	ctx := t.Context()
	ch, err := handler.Start(ctx, cfg)
	require.NoError(t, err)

	select {
	case evt := <-ch:
		require.Equal(t, wsuser.EventOrder, evt.EventType)
		require.Equal(t, wsuser.OrderCancellation, evt.OrderType)
		require.Equal(t, "order-abc", evt.OrderID)
		require.Equal(t, "0xdeadbeef", evt.Market)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for user event")
	}
}
