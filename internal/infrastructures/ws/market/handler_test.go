package market_test

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

	wsmarket "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/market"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func TestMarketHandler_Start_ReceivesTickSizeChange(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// read subscription message
		_, _, _ = conn.ReadMessage()

		msg, _ := json.Marshal(map[string]any{
			"event_type": "tick_size_change",
			"data": map[string]string{
				"asset_id": "0xcondition123",
				"new_size": "0.001",
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
	handler := wsmarket.New(wsURL)

	ctx := t.Context()
	ch, err := handler.Start(ctx, []string{"0xcondition123"})
	require.NoError(t, err)

	select {
	case evt := <-ch:
		assert.Equal(t, wsmarket.EventTickSizeChange, evt.Type)
		require.NotNil(t, evt.TickSizeChange)
		assert.Equal(t, "0xcondition123", evt.TickSizeChange.ConditionID)
		assert.Equal(t, "0.001", evt.TickSizeChange.NewTickSize)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for market event")
	}
}
