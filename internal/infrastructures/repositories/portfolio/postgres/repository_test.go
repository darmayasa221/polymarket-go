package postgres_test

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
	pospostgres "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/portfolio/postgres"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
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

func makePosition(t *testing.T) *position.Position {
	t.Helper()
	pos, err := position.New(position.Params{
		Asset:    market.BTC,
		TokenID:  polyid.TokenID("token-up"),
		Outcome:  market.Up,
		Size:     decimal.RequireFromString("10"),
		AvgPrice: decimal.RequireFromString("0.60"),
		MarketID: "market-1",
	})
	require.NoError(t, err)
	return pos
}

func TestPositionRepository_SaveAndFindByID(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pospostgres.New(newTestDB(t))

	pos := makePosition(t)
	require.NoError(t, repo.Save(ctx, pos))

	got, err := repo.FindByID(ctx, pos.ID())
	require.NoError(t, err)
	assert.Equal(t, pos.ID(), got.ID())
	assert.Equal(t, market.BTC, got.Asset())
	assert.Nil(t, got.ClosedAt())
}

func TestPositionRepository_ListOpen(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pospostgres.New(newTestDB(t))

	require.NoError(t, repo.Save(ctx, makePosition(t)))
	require.NoError(t, repo.Save(ctx, makePosition(t)))

	list, err := repo.ListOpen(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestPositionRepository_Close(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pospostgres.New(newTestDB(t))

	pos := makePosition(t)
	require.NoError(t, repo.Save(ctx, pos))

	exitPrice := decimal.RequireFromString("0.80")
	closedAt := timeutil.Now()
	require.NoError(t, repo.Close(ctx, pos.ID(), exitPrice, closedAt))

	open, err := repo.ListOpen(ctx)
	require.NoError(t, err)
	assert.Empty(t, open)
}

func TestPositionRepository_ListClosedWithExitPrice(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pospostgres.New(newTestDB(t))

	pos := makePosition(t)
	require.NoError(t, repo.Save(ctx, pos))

	exitPrice := decimal.RequireFromString("0.80")
	closedAt := timeutil.Now()
	require.NoError(t, repo.Close(ctx, pos.ID(), exitPrice, closedAt))

	records, err := repo.ListClosedWithExitPrice(ctx)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.True(t, records[0].ExitPrice.Equal(exitPrice))
}

func TestPositionRepository_FindByID_NotFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := pospostgres.New(newTestDB(t))
	_, err := repo.FindByID(ctx, "nonexistent")
	require.Error(t, err)
}
