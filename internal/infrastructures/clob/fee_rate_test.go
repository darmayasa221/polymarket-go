package clob_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
)

func TestFeeRateProvider_FetchFeeRate(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/fee-rate", r.URL.Path)
		assert.Equal(t, "token-123", r.URL.Query().Get("token_id"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{"base_fee": 1000})
	}))
	t.Cleanup(srv.Close)

	cfg := clob.Config{BaseURL: srv.URL, APISecret: validTestSecret()}
	provider := clob.NewFeeRateProvider(clob.NewClient(cfg))

	fee, err := provider.FetchFeeRate(t.Context(), "token-123")
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), fee)
}
