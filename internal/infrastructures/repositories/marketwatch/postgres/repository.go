// Package postgres implements the MarketRepository using PostgreSQL.
package postgres

import (
	"database/sql"

	marketwatchports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
)

// Compile-time assertion: Repository implements marketwatchports.MarketRepository.
var _ marketwatchports.MarketRepository = (*Repository)(nil)

// Repository is the PostgreSQL implementation of marketwatchports.MarketRepository.
type Repository struct{ db *sql.DB }

// New creates a new PostgreSQL market repository.
func New(db *pgdb.DB) *Repository { return &Repository{db: db.DB()} }
