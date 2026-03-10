package postgres_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
)

func TestRunMigrations(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run postgres integration tests")
	}
	t.Parallel()

	db, err := pgdb.New(pgdb.Config{DSN: dsn})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	// First call — must succeed.
	require.NoError(t, pgdb.RunMigrations(db.DB()))

	// Second call — idempotent (CREATE TABLE IF NOT EXISTS).
	require.NoError(t, pgdb.RunMigrations(db.DB()))
}
