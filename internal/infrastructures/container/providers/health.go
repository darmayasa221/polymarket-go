package providers

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/health"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
	infrahealth "github.com/darmayasa221/polymarket-go/internal/infrastructures/health"
)

// ProvideHealthChecker wires all health components into a Checker.
func ProvideHealthChecker(db *sqlitedb.DB, cache *redisclient.Client) *health.Checker {
	return health.NewChecker(
		infrahealth.NewDatabaseChecker(db),
		infrahealth.NewCacheChecker(cache),
	)
}
