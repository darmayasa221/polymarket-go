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

func TestUserHandler_Start_ReceivesOrderFilled(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// read subscription message
		_, _, _ = conn.ReadMessage()

		msg, _ := json.Marshal(map[string]any{
			"event_type": "order_filled",
			"data": map[string]string{
				"id": "order-abc",
			},
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
		require.Equal(t, wsuser.EventOrderFilled, evt.Type)
		require.Equal(t, "order-abc", evt.OrderID)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for user event")
	}
}
