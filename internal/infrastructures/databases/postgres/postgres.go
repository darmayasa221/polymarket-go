package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver via database/sql.
)

// DB wraps the standard sql.DB for PostgreSQL.
type DB struct {
	db *sql.DB
}

// New opens a PostgreSQL connection using the given config.
func New(cfg Config) (*DB, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return &DB{db: db}, nil
}

// DB returns the underlying *sql.DB.
func (p *DB) DB() *sql.DB { return p.db }

// Close closes the database connection.
func (p *DB) Close() error { return p.db.Close() }

// Ping verifies the connection is still alive.
func (p *DB) Ping(ctx context.Context) error { return p.db.PingContext(ctx) }
