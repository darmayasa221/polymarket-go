package clob_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
)

func TestHeartbeatSender_Send(t *testing.T) {
	t.Parallel()

	var callCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/heartbeats", r.URL.Path)

		n := callCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			// First call: seedChain sends null heartbeat_id.
			// Server returns 400 with seed ID embedded in error response.
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"heartbeat_id": "seed-001"})
			return
		}
		// Second call: valid chain heartbeat — return 200 with next ID.
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"heartbeat_id": "next-001"})
	}))
	t.Cleanup(srv.Close)

	cfg := clob.Config{BaseURL: srv.URL, APISecret: validTestSecret()}
	sender := clob.NewHeartbeatSender(clob.NewClient(cfg))
	require.NoError(t, sender.Send(t.Context()))
	assert.EqualValues(t, 2, callCount.Load(), "expected 2 HTTP calls: seed + chain")
}
