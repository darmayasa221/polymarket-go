// Package health provides infrastructure health check implementations.
package health

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/health"
	healthtypes "github.com/darmayasa221/polymarket-go/internal/applications/health/types"
	healthconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/health"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// Compile-time assertion: DatabaseChecker implements health.Component.
var _ health.Component = (*DatabaseChecker)(nil)

// DatabaseChecker checks SQLite connectivity.
type DatabaseChecker struct {
	db *sqlitedb.DB
}

// NewDatabaseChecker creates a new database health checker.
func NewDatabaseChecker(db *sqlitedb.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

// Name returns the component name.
func (c *DatabaseChecker) Name() string { return "database" }

// Check verifies the database is reachable.
func (c *DatabaseChecker) Check(ctx context.Context) healthtypes.ComponentStatus {
	if err := c.db.Ping(ctx); err != nil {
		return healthtypes.ComponentStatus{
			Name:    c.Name(),
			Status:  healthconst.StatusUnhealthy,
			Message: err.Error(),
		}
	}
	return healthtypes.ComponentStatus{Name: c.Name(), Status: healthconst.StatusHealthy}
}
