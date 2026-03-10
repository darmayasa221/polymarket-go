// Package config provides the root configuration facade that aggregates all component configs.
package config

import (
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
	httpserver "github.com/darmayasa221/polymarket-go/internal/infrastructures/interfaces/http/server"
	bcrypt "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/encryption/bcrypt"
	jwt "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/token/jwt"
)

// AppConfig holds general application-level configuration.
type AppConfig struct {
	// Name is the application name.
	Name string
	// Environment is the runtime environment (e.g. "development", "production").
	Environment string
	// Version is the application version string.
	Version string
	// LogLevel is the minimum log level to emit (e.g. "info", "debug").
	LogLevel string
}

// PostgreSQLConfig holds the PostgreSQL database configuration.
type PostgreSQLConfig struct {
	// DSN is the PostgreSQL data source name.
	// Example: "postgres://user:pass@localhost:5432/polymarket?sslmode=disable"
	DSN string
}

// CLOBConfig holds Polymarket CLOB API credentials.
type CLOBConfig struct {
	// BaseURL is the CLOB REST base URL.
	BaseURL string
	// Address is the EOA wallet address (funder).
	Address string
	// APIKey is the derived Polymarket API key.
	APIKey string
	// APISecret is the base64-encoded HMAC secret.
	APISecret string
	// APIPassphrase is the API passphrase.
	APIPassphrase string
}

// GammaConfig holds Polymarket Gamma API configuration.
type GammaConfig struct {
	// BaseURL is the Gamma API base URL.
	BaseURL string
}

// Config is the root configuration facade that aggregates all component configurations.
type Config struct {
	// App holds general application configuration.
	App AppConfig
	// HTTP holds the HTTP server configuration.
	HTTP httpserver.Config
	// Database holds the SQLite database configuration.
	Database sqlite.Config
	// PostgreSQL holds the PostgreSQL database configuration.
	PostgreSQL PostgreSQLConfig
	// Cache holds the Redis cache configuration.
	Cache redis.Config
	// JWT holds the JWT token manager configuration.
	JWT jwt.Config
	// Bcrypt holds the bcrypt encryption configuration.
	Bcrypt bcrypt.Config
	// CLOB holds the Polymarket CLOB API configuration.
	CLOB CLOBConfig
	// Gamma holds the Polymarket Gamma API configuration.
	Gamma GammaConfig
}
