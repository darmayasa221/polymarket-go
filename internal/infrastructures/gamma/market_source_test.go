package gamma_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/gamma"
)

func TestMarketSource_FetchActive5mMarkets(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/events", r.URL.Path)
		slugs := r.URL.Query()["slug"]
		require.NotEmpty(t, slugs, "expected slug query params")

		// Find the BTC slug from the request to use in the response.
		btcSlug := ""
		for _, s := range slugs {
			if strings.HasPrefix(s, "btc-updown-5m-") {
				btcSlug = s
				break
			}
		}
		require.NotEmpty(t, btcSlug, "btc slug not found in query")

		// Gamma API returns []gammaEvent; each event has a nested markets[] array.
		// outcomes and clobTokenIds are JSON-encoded strings (double-encoded).
		payload := []map[string]any{
			{
				"id":     "event-1",
				"slug":   btcSlug,
				"closed": false,
				"markets": []map[string]any{
					{
						"id":                    "inner-1",
						"conditionId":           "0xabc123",
						"slug":                  btcSlug,
						"outcomes":              `["Up","Down"]`,
						"clobTokenIds":          `["tok-up-1","tok-down-1"]`,
						"orderPriceMinTickSize": 0.01,
						"enableOrderBook":       true,
						"active":                true,
						"closed":                false,
					},
				},
			},
			{
				// closed event — should be excluded
				"id":      "event-2",
				"slug":    "eth-updown-5m-0",
				"closed":  true,
				"markets": []map[string]any{},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	t.Cleanup(srv.Close)

	src := gamma.NewMarketSource(gamma.Config{BaseURL: srv.URL})
	markets, err := src.FetchActive5mMarkets(t.Context())
	require.NoError(t, err)

	// Only the open BTC market should be returned.
	require.Len(t, markets, 1)
	m := markets[0]
	assert.Equal(t, "btc", string(m.Asset()))
	assert.Equal(t, "0xabc123", string(m.ConditionID()))
	assert.Equal(t, "tok-up-1", string(m.UpTokenID()))
	assert.Equal(t, "tok-down-1", string(m.DownTokenID()))
}
