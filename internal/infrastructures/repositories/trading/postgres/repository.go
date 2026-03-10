// Package postgres implements the OrderRepository using PostgreSQL.
package postgres

import (
	"database/sql"

	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
)

// Compile-time assertion: Repository implements tradingports.OrderRepository.
var _ tradingports.OrderRepository = (*Repository)(nil)

// Repository is the PostgreSQL implementation of tradingports.OrderRepository.
type Repository struct{ db *sql.DB }

// New creates a new PostgreSQL order repository.
func New(db *pgdb.DB) *Repository { return &Repository{db: db.DB()} }
