// Package sqlite implements the Authentication repository using SQLite (Adapter Pattern).
// This is an alternative persistent implementation of repository.Authentication.
// The production DI container currently wires the Redis implementation for performance.
// This package is available as a drop-in replacement for environments without Redis
// or as a durable fallback store alongside the Redis adapter.
package sqlite

import (
	"database/sql"

	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// Compile-time assertion: Repository implements repository.Authentication.
var _ repository.Authentication = (*Repository)(nil)

// Repository is the SQLite implementation of repository.Authentication.
type Repository struct {
	db *sql.DB
}

// New creates a new SQLite Authentication repository.
func New(db *sqlitedb.DB) *Repository {
	return &Repository{db: db.DB()}
}
