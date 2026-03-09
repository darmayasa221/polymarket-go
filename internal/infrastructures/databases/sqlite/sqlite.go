// Package sqlite provides the SQLite database adapter.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // CGO-free SQLite driver.
)

// DB wraps the standard sql.DB for SQLite.
type DB struct {
	db *sql.DB
}

// New opens a SQLite database connection using the given config.
// The startup ping uses context.Background — callers must impose their own deadline if needed.
func New(cfg Config) (*DB, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("sqlite: ping: %w", err)
	}
	// SQLite best practices for concurrent access.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return &DB{db: db}, nil
}

// DB returns the underlying *sql.DB for direct query use.
func (s *DB) DB() *sql.DB { return s.db }

// Close closes the database connection.
func (s *DB) Close() error { return s.db.Close() }

// Ping verifies the connection is still alive.
func (s *DB) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
