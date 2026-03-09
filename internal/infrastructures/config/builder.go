package config

import (
	appconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/app"
	cacheconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/cache"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
	httpserver "github.com/darmayasa221/polymarket-go/internal/infrastructures/interfaces/http/server"
	bcrypt "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/encryption/bcrypt"
	jwt "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/token/jwt"
)

// Test-only constants used by NewBuilder. These are intentionally not in commons
// because they are specific to the test config builder and have no production use.
const (
	testJWTSecretKey = "test-secret-key-minimum-32-chars!!"
	testMaxBodyBytes = 1 << 20 // 1 MiB — intentionally smaller than production 10 MiB
)

// Builder constructs a Config programmatically, primarily for use in tests.
// It is initialized with safe test defaults so callers only override what they need.
type Builder struct {
	cfg *Config
}

// NewBuilder returns a Builder populated with safe test defaults:
//   - SQLite in-memory database
//   - bcrypt cost 4 (minimum, fast for tests)
//   - Redis at localhost:6379
//   - JWT access 15 m / refresh 168 h with a non-empty secret
func NewBuilder() *Builder {
	return &Builder{
		cfg: &Config{
			App: AppConfig{
				Name:        appconst.Name,
				Environment: appconst.EnvTest,
				Version:     appconst.TestVersion,
				LogLevel:    appconst.LogLevelDebug,
			},
			HTTP: httpserver.Config{
				Port:              appconst.DefaultHTTPPort,
				ReadTimeout:       appconst.DefaultReadTimeoutDuration,
				WriteTimeout:      appconst.DefaultWriteTimeoutDuration,
				IdleTimeout:       appconst.DefaultIdleTimeoutDuration,
				RequestTimeout:    appconst.DefaultTimeoutDuration,
				MaxBodyBytes:      testMaxBodyBytes,
				AllowedOrigins:    []string{appconst.DefaultAllowedOrigins},
				RateLimitRequests: appconst.DefaultRateLimitRequests,
				RateLimitWindow:   appconst.DefaultRateLimitWindow,
			},
			Database: sqlite.Config{
				Path: ":memory:",
			},
			Cache: redis.Config{
				Address:     cacheconst.DefaultAddress,
				DB:          0,
				DialTimeout: cacheconst.DefaultDialTimeoutDuration,
				ReadTimeout: cacheconst.DefaultReadTimeoutDuration,
			},
			JWT: jwt.Config{
				SecretKey:            testJWTSecretKey,
				AccessTokenDuration:  tokenconst.DefaultAccessTokenDuration,
				RefreshTokenDuration: tokenconst.DefaultRefreshTokenDuration,
				Issuer:               appconst.Name,
			},
			Bcrypt: bcrypt.Config{
				Cost: 4,
			},
		},
	}
}

// WithDatabase overrides the SQLite database path.
func (b *Builder) WithDatabase(path string) *Builder {
	b.cfg.Database.Path = path
	return b
}

// WithCache overrides the Redis server address.
func (b *Builder) WithCache(address string) *Builder {
	b.cfg.Cache.Address = address
	return b
}

// Build returns the constructed Config.
func (b *Builder) Build() *Config {
	c := *b.cfg
	return &c
}
