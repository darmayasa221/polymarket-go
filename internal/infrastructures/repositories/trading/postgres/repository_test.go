package postgres_test

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
	orderpostgres "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/trading/postgres"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
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

func makeOrder(t *testing.T) *order.Order {
	t.Helper()
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
	return o
}

func TestOrderRepository_SaveAndFindByID(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := orderpostgres.New(newTestDB(t))

	o := makeOrder(t)
	require.NoError(t, repo.Save(ctx, o))

	got, err := repo.FindByID(ctx, o.ID())
	require.NoError(t, err)
	assert.Equal(t, o.ID(), got.ID())
	assert.Equal(t, order.StatusOpen, got.Status())
	assert.Equal(t, order.Buy, got.Side())
}

func TestOrderRepository_ListOpenByMarket(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := orderpostgres.New(newTestDB(t))

	require.NoError(t, repo.Save(ctx, makeOrder(t)))
	require.NoError(t, repo.Save(ctx, makeOrder(t)))

	list, err := repo.ListOpenByMarket(ctx, "market-1")
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := orderpostgres.New(newTestDB(t))

	o := makeOrder(t)
	require.NoError(t, repo.Save(ctx, o))
	require.NoError(t, repo.UpdateStatus(ctx, o.ID(), order.StatusFilled))

	got, err := repo.FindByID(ctx, o.ID())
	require.NoError(t, err)
	assert.Equal(t, order.StatusFilled, got.Status())
}

func TestOrderRepository_FindByID_NotFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := orderpostgres.New(newTestDB(t))
	_, err := repo.FindByID(ctx, polyid.OrderID("nonexistent"))
	require.Error(t, err)
}

func makeGTDOrder(t *testing.T) *order.Order {
	t.Helper()
	o, err := order.New(order.Params{
		MarketID:   "market-1",
		TokenID:    polyid.TokenID("token-up"),
		Side:       order.Buy,
		Outcome:    market.Up,
		Price:      decimal.RequireFromString("0.60"),
		Size:       decimal.RequireFromString("10"),
		Type:       order.GTD,
		Expiration: time.Now().UTC().Add(5 * time.Minute),
		FeeRateBps: 100,
	})
	require.NoError(t, err)
	return o
}

func TestOrderRepository_GTDWithExpiration(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := orderpostgres.New(newTestDB(t))

	o := makeGTDOrder(t)
	require.NoError(t, repo.Save(ctx, o))

	got, err := repo.FindByID(ctx, o.ID())
	require.NoError(t, err)
	assert.Equal(t, order.GTD, got.Type())
	assert.False(t, got.Expiration().IsZero())
}
