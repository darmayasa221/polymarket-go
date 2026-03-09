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

	if len(missing) > 0 {
		return errors.New("missing required configuration:\n" + strings.Join(missing, "\n"))
	}

	return nil
}
