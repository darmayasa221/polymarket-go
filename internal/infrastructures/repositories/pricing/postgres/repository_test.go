package postgres_test

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
	pricepostgres "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/pricing/postgres"
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

func TestPriceRepository_SaveAndLatestByAsset(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pricepostgres.New(newTestDB(t))

	p, err := oracle.New(oracle.Params{
		Asset:      "btc",
		Source:     oracle.SourceBinance,
		Value:      decimal.RequireFromString("65000.00"),
		ReceivedAt: time.Now().UTC(),
	})
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, p))

	got, err := repo.LatestByAsset(ctx, "btc")
	require.NoError(t, err)
	assert.Equal(t, "btc", got.Asset())
	assert.Equal(t, oracle.SourceBinance, got.Source())
	assert.True(t, got.Value().Equal(p.Value()))
}

func TestPriceRepository_LatestChainlinkByAsset(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pricepostgres.New(newTestDB(t))

	// Binance should be ignored by LatestChainlink.
	binance, _ := oracle.New(oracle.Params{
		Asset: "eth", Source: oracle.SourceBinance,
		Value: decimal.RequireFromString("3000"), ReceivedAt: time.Now().UTC(),
	})
	require.NoError(t, repo.Save(ctx, binance))

	chainlink, _ := oracle.New(oracle.Params{
		Asset: "eth", Source: oracle.SourceChainlink,
		Value: decimal.RequireFromString("3001"), ReceivedAt: time.Now().UTC(),
	})
	require.NoError(t, repo.Save(ctx, chainlink))

	got, err := repo.LatestChainlinkByAsset(ctx, "eth")
	require.NoError(t, err)
	assert.Equal(t, oracle.SourceChainlink, got.Source())
}

func TestPriceRepository_WindowOpenPrice(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pricepostgres.New(newTestDB(t))

	windowStart := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)

	// Before window — must NOT be returned.
	before, _ := oracle.New(oracle.Params{
		Asset: "btc", Source: oracle.SourceChainlink,
		Value: decimal.RequireFromString("64000"), ReceivedAt: windowStart.Add(-time.Second),
	})
	require.NoError(t, repo.Save(ctx, before))

	// At window start — must be returned.
	atStart, _ := oracle.New(oracle.Params{
		Asset: "btc", Source: oracle.SourceChainlink,
		Value: decimal.RequireFromString("65000"), ReceivedAt: windowStart,
	})
	require.NoError(t, repo.Save(ctx, atStart))

	got, err := repo.WindowOpenPrice(ctx, "btc", windowStart)
	require.NoError(t, err)
	assert.True(t, got.Value().Equal(decimal.RequireFromString("65000")))
}

func TestPriceRepository_LatestByAsset_NotFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pricepostgres.New(newTestDB(t))
	_, err := repo.LatestByAsset(ctx, "sol")
	require.Error(t, err)
}
