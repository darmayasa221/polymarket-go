// Package postgres implements the PositionRepository using PostgreSQL.
package postgres

import (
	"database/sql"

	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
)

// Compile-time assertion: Repository implements portfolioports.PositionRepository.
var _ portfolioports.PositionRepository = (*Repository)(nil)

// Repository is the PostgreSQL implementation of portfolioports.PositionRepository.
type Repository struct{ db *sql.DB }

// New creates a new PostgreSQL position repository.
func New(db *pgdb.DB) *Repository { return &Repository{db: db.DB()} }
