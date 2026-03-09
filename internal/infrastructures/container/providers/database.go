// Package providers wires infrastructure adapters into the Container.
package providers

import (
	"fmt"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// ProvideDatabase creates and returns a SQLite database connection.
func ProvideDatabase(cfg *config.Config) (*sqlitedb.DB, error) {
	db, err := sqlitedb.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("providers: database: %w", err)
	}
	return db, nil
}
