package config

import (
	"errors"
	"strings"
)

// validate checks that all required configuration fields are present.
// It returns an error listing every missing field if any are absent.
func validate(cfg *Config) error {
	var missing []string

	if cfg.JWT.SecretKey == "" {
		missing = append(missing, "  - JWT_SECRET_KEY (jwt.secret_key)")
	}

	if cfg.Database.Path == "" {
		missing = append(missing, "  - DB_PATH (database.path)")
	}

	if cfg.Cache.Address == "" {
		missing = append(missing, "  - REDIS_ADDRESS (cache.address)")
	}

	if cfg.PostgreSQL.DSN == "" {
		missing = append(missing, "  - DATABASE_URL (postgresql.dsn)")
	}

	if cfg.CLOB.APIKey == "" {
		missing = append(missing, "  - POLY_API_KEY (clob.api_key)")
	}

	if cfg.CLOB.APISecret == "" {
		missing = append(missing, "  - POLY_API_SECRET (clob.api_secret)")
	}

	if cfg.CLOB.Address == "" {
		missing = append(missing, "  - POLY_FUNDER_ADDRESS (clob.address)")
	}

	if len(missing) > 0 {
		return errors.New("missing required configuration:\n" + strings.Join(missing, "\n"))
	}

	return nil
}
