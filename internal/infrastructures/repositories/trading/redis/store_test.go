package redis_test

import (
	"os"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/applications/shared/windowstate"
	redisstore "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/trading/redis"
)

func newTestClient(t *testing.T) *goredis.Client {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("set TEST_REDIS_URL to run redis integration tests")
	}
	opts, err := goredis.ParseURL(url)
	require.NoError(t, err)
	client := goredis.NewClient(opts)
	t.Cleanup(func() {
		_ = client.FlushDB(t.Context())
		_ = client.Close()
	})
	return client
}

func TestWindowStateStore_SaveAndGet(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	store := redisstore.New(newTestClient(t))

	state := windowstate.WindowState{
		Asset:       "btc",
		WindowStart: time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
		Status:      windowstate.WindowOpen,
		OpenPrice:   decimal.RequireFromString("0.60"),
	}
	require.NoError(t, store.SaveWindowState(ctx, state))

	got, err := store.GetWindowState(ctx, "btc")
	require.NoError(t, err)
	assert.Equal(t, state.Asset, got.Asset)
	assert.Equal(t, windowstate.WindowOpen, got.Status)
	assert.True(t, got.OpenPrice.Equal(state.OpenPrice))
}

func TestWindowStateStore_GetNotFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	store := redisstore.New(newTestClient(t))
	_, err := store.GetWindowState(ctx, "sol")
	require.Error(t, err)
}

func TestWindowStateStore_Overwrite(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	store := redisstore.New(newTestClient(t))

	state1 := windowstate.WindowState{Asset: "eth", Status: windowstate.WindowOpen, OpenPrice: decimal.RequireFromString("0.60")}
	state2 := windowstate.WindowState{Asset: "eth", Status: windowstate.WindowClosed, OpenPrice: decimal.RequireFromString("0.80")}

	require.NoError(t, store.SaveWindowState(ctx, state1))
	require.NoError(t, store.SaveWindowState(ctx, state2))

	got, err := store.GetWindowState(ctx, "eth")
	require.NoError(t, err)
	assert.Equal(t, windowstate.WindowClosed, got.Status)
	assert.True(t, got.OpenPrice.Equal(decimal.RequireFromString("0.80")))
}
