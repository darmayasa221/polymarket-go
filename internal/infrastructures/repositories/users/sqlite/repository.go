// Package sqlite implements the User repository using SQLite (Adapter Pattern).
package sqlite

import (
	"database/sql"

	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
)

// Compile-time assertion: Repository implements repository.User.
var _ repository.User = (*Repository)(nil)

// Repository is the SQLite implementation of repository.User.
type Repository struct {
	db *sql.DB
}

// New creates a new SQLite User repository.
func New(db *sqlitedb.DB) *Repository {
	return &Repository{db: db.DB()}
}
