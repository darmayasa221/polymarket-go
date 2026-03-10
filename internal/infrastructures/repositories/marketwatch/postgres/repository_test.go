package postgres_test

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
	marketpostgres "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/marketwatch/postgres"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

func newTestDB(t *testing.T) *pgdb.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run postgres integration tests")
	}
	db, err := pgdb.New(pgdb.Config{DSN: dsn})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, pgdb.RunMigrations(db.DB()))
	_, _ = db.DB().ExecContext(t.Context(), "TRUNCATE prices, markets, orders, positions RESTART IDENTITY CASCADE")
	return db
}

func makeMarket(t *testing.T, asset market.Asset, windowStart time.Time) *market.Market {
	t.Helper()
	m, err := market.New(market.Params{
		ID:          "gamma-event-1",
		Asset:       asset,
		WindowStart: windowStart,
		ConditionID: polyid.ConditionID("cond-1"),
		UpTokenID:   polyid.TokenID("tok-up"),
		DownTokenID: polyid.TokenID("tok-down"),
		TickSize:    decimal.NewFromFloat(0.01),
		FeeEnabled:  true,
	})
	require.NoError(t, err)
	return m
}

func TestMarketRepository_SaveAndFindByAssetAndWindow(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := marketpostgres.New(newTestDB(t))

	windowStart := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	m := makeMarket(t, market.BTC, windowStart)

	require.NoError(t, repo.Save(ctx, m))

	got, err := repo.FindByAssetAndWindow(ctx, market.BTC, windowStart)
	require.NoError(t, err)
	assert.Equal(t, m.ID(), got.ID())
	assert.Equal(t, market.BTC, got.Asset())
	assert.True(t, got.Active())
}

func TestMarketRepository_UpdateTickSize(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := marketpostgres.New(newTestDB(t))

	windowStart := time.Date(2026, 3, 10, 12, 5, 0, 0, time.UTC)
	m := makeMarket(t, market.ETH, windowStart)
	require.NoError(t, repo.Save(ctx, m))

	newTick := decimal.NewFromFloat(0.001)
	require.NoError(t, repo.UpdateTickSize(ctx, string(m.ConditionID()), newTick))

	got, err := repo.FindByAssetAndWindow(ctx, market.ETH, windowStart)
	require.NoError(t, err)
	assert.True(t, got.TickSize().Equal(newTick))
}

func TestMarketRepository_ListActive(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := marketpostgres.New(newTestDB(t))

	w1 := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	w2 := time.Date(2026, 3, 10, 12, 5, 0, 0, time.UTC)

	require.NoError(t, repo.Save(ctx, makeMarket(t, market.SOL, w1)))
	require.NoError(t, repo.Save(ctx, makeMarket(t, market.XRP, w2)))

	list, err := repo.ListActive(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestMarketRepository_FindByAssetAndWindow_NotFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := marketpostgres.New(newTestDB(t))
	windowStart := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	_, err := repo.FindByAssetAndWindow(ctx, market.BTC, windowStart)
	require.Error(t, err)
}
