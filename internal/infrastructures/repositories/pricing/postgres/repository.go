// Package postgres implements the PriceRepository using PostgreSQL.
package postgres

import (
	"database/sql"

	pricingports "github.com/darmayasa221/polymarket-go/internal/applications/pricing/ports"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
)

// Compile-time assertion: Repository implements pricingports.PriceRepository.
var _ pricingports.PriceRepository = (*Repository)(nil)

// Repository is the PostgreSQL implementation of pricingports.PriceRepository.
type Repository struct {
	db *sql.DB
}

// New creates a new PostgreSQL price repository.
func New(db *pgdb.DB) *Repository {
	return &Repository{db: db.DB()}
}
