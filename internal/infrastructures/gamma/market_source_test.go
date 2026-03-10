package gamma_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/gamma"
)

func TestMarketSource_FetchActive5mMarkets(t *testing.T) {
	t.Parallel()

	payload := []map[string]any{
		{
			"id":          "event-1",
			"slug":        "btc-updown-5m-1741996800",
			"conditionId": "0xabc123",
			"tokens": []map[string]string{
				{"outcome": "Up", "token_id": "tok-up-1"},
				{"outcome": "Down", "token_id": "tok-down-1"},
			},
			"minimum_tick_size": "0.01",
			"fees_enabled":      true,
			"enable_order_book": true,
			"closed":            false,
		},
		{
			// closed market — should be excluded
			"id":          "event-2",
			"slug":        "eth-updown-5m-1741996800",
			"conditionId": "0xdef456",
			"tokens": []map[string]string{
				{"outcome": "Up", "token_id": "tok-up-2"},
				{"outcome": "Down", "token_id": "tok-down-2"},
			},
			"minimum_tick_size": "0.01",
			"fees_enabled":      true,
			"enable_order_book": true,
			"closed":            true,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/events", r.URL.Path)
		assert.Equal(t, "102892", r.URL.Query().Get("tag"))
		assert.Equal(t, "false", r.URL.Query().Get("closed"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	t.Cleanup(srv.Close)

	src := gamma.NewMarketSource(gamma.Config{BaseURL: srv.URL})
	markets, err := src.FetchActive5mMarkets(t.Context())
	require.NoError(t, err)

	// only the open market should be returned
	require.Len(t, markets, 1)
	m := markets[0]
	assert.Equal(t, "btc", string(m.Asset()))
	assert.Equal(t, "0xabc123", string(m.ConditionID()))
	assert.Equal(t, "tok-up-1", string(m.UpTokenID()))
	assert.Equal(t, "tok-down-1", string(m.DownTokenID()))
}
