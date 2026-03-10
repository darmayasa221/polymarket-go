package clob_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
)

func TestOrderSubmitter_Submit(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/order", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"orderID": "clob-order-abc"})
	}))
	t.Cleanup(srv.Close)

	o, err := order.New(order.Params{
		MarketID:   "market-1",
		TokenID:    polyid.TokenID("token-up"),
		Side:       order.Buy,
		Outcome:    market.Up,
		Price:      decimal.RequireFromString("0.60"),
		Size:       decimal.RequireFromString("10"),
		Type:       order.FOK,
		FeeRateBps: 100,
	})
	require.NoError(t, err)

	cfg := clob.Config{BaseURL: srv.URL, APISecret: validTestSecret()}
	submitter := clob.NewOrderSubmitter(clob.NewClient(cfg))

	clobID, err := submitter.Submit(t.Context(), o, []byte("signature-bytes"))
	require.NoError(t, err)
	assert.Equal(t, "clob-order-abc", clobID)
}

func TestOrderSubmitter_Cancel(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	cfg := clob.Config{BaseURL: srv.URL, APISecret: validTestSecret()}
	require.NoError(t, clob.NewOrderSubmitter(clob.NewClient(cfg)).Cancel(t.Context(), "order-abc"))
}
