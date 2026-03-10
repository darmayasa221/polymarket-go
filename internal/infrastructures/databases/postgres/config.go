// Package postgres provides the PostgreSQL database adapter.
package postgres

// Config holds the configuration for the PostgreSQL database adapter.
type Config struct {
	// DSN is the PostgreSQL data source name.
	// Example: "postgres://user:pass@localhost:5432/polymarket?sslmode=disable"
	DSN string
}
