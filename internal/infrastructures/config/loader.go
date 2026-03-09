package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	appconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/app"
	cacheconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/cache"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
	httpserver "github.com/darmayasa221/polymarket-go/internal/infrastructures/interfaces/http/server"
	bcrypt "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/encryption/bcrypt"
	jwt "github.com/darmayasa221/polymarket-go/internal/infrastructures/security/token/jwt"
)

// Load reads configuration from the given env file and environment variables.
// If the env file is missing, it is silently ignored. Returns an error if
// required configuration values are absent.
func Load(envFile string) (*Config, error) {
	_ = godotenv.Load(envFile)
	viper.AutomaticEnv()

	cfg := &Config{
		App:      appConfig(),
		HTTP:     httpServerConfig(),
		Database: sqliteConfig(),
		Cache:    redisConfig(),
		JWT:      jwtConfig(),
		Bcrypt:   bcryptConfig(),
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// appConfig builds the AppConfig section from environment variables.
func appConfig() AppConfig {
	return AppConfig{
		Name:        getOrDefault("APP_NAME", appconst.Name),
		Environment: getOrDefault("APP_ENV", appconst.EnvDevelopment),
		Version:     getOrDefault("APP_VERSION", appconst.Version),
		LogLevel:    getOrDefault("APP_LOG_LEVEL", appconst.DefaultLogLevel),
	}
}

// sqliteConfig builds the sqlite.Config section from environment variables.
func sqliteConfig() sqlite.Config {
	return sqlite.Config{
		Path: viper.GetString("DB_PATH"),
	}
}

// redisConfig builds the redis.Config section from environment variables.
func redisConfig() redis.Config {
	return redis.Config{
		Address:     viper.GetString("REDIS_ADDRESS"),
		Password:    viper.GetString("REDIS_PASSWORD"),
		DB:          viper.GetInt("REDIS_DB"),
		DialTimeout: parseDuration("REDIS_DIAL_TIMEOUT", cacheconst.DefaultDialTimeout),
		ReadTimeout: parseDuration("REDIS_READ_TIMEOUT", cacheconst.DefaultReadTimeout),
	}
}

// jwtConfig builds the jwt.Config section from environment variables.
func jwtConfig() jwt.Config {
	return jwt.Config{
		SecretKey:            viper.GetString("JWT_SECRET_KEY"),
		AccessTokenDuration:  parseDuration("JWT_ACCESS_TOKEN_DURATION", tokenconst.DefaultAccessTokenDurationStr),
		RefreshTokenDuration: parseDuration("JWT_REFRESH_TOKEN_DURATION", tokenconst.DefaultRefreshTokenDurationStr),
		Issuer:               getOrDefault("JWT_ISSUER", appconst.Name),
	}
}

// bcryptConfig builds the bcrypt.Config section from environment variables.
// Zero cost is passed as-is; bcrypt.New applies the DefaultCost guard.
func bcryptConfig() bcrypt.Config {
	return bcrypt.Config{Cost: viper.GetInt("BCRYPT_COST")}
}

// httpServerConfig builds the httpserver.Config section from environment variables.
func httpServerConfig() httpserver.Config {
	return httpserver.Config{
		Port:              getOrDefault("HTTP_PORT", appconst.DefaultHTTPPort),
		ReadTimeout:       parseDuration("HTTP_READ_TIMEOUT", appconst.DefaultReadTimeout),
		WriteTimeout:      parseDuration("HTTP_WRITE_TIMEOUT", appconst.DefaultWriteTimeout),
		IdleTimeout:       parseDuration("HTTP_IDLE_TIMEOUT", appconst.DefaultIdleTimeout),
		RequestTimeout:    parseDuration("HTTP_REQUEST_TIMEOUT", appconst.DefaultTimeout),
		MaxBodyBytes:      parseMaxBodyBytes("HTTP_MAX_BODY_BYTES"),
		AllowedOrigins:    parseCORSOrigins("CORS_ALLOWED_ORIGINS"),
		RateLimitRequests: getIntOrDefault("RATE_LIMIT_REQUESTS", appconst.DefaultRateLimitRequests),
		RateLimitWindow:   parseDuration("RATE_LIMIT_WINDOW", appconst.DefaultRateLimitWindowStr),
	}
}

// maxInt64AsUint is the maximum int64 value expressed as uint for overflow checks.
const maxInt64AsUint = uint(1<<63 - 1)

// parseMaxBodyBytes reads the maximum body size from viper by key and returns
// it as int64. When the key is unset or zero, it returns 0 (no limit applied by default).
func parseMaxBodyBytes(key string) int64 {
	v := viper.GetSizeInBytes(key)
	if v > maxInt64AsUint {
		return int64(maxInt64AsUint)
	}

	return int64(v)
}

// parseDuration reads a duration string from viper by key and falls back to
// defaultVal if the key is unset or unparseable.
func parseDuration(key, defaultVal string) time.Duration {
	raw := viper.GetString(key)
	if raw == "" {
		d, _ := time.ParseDuration(defaultVal)
		return d
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		d, _ = time.ParseDuration(defaultVal)
	}

	return d
}

// getOrDefault returns the viper string value for key, or defaultVal when empty.
func getOrDefault(key, defaultVal string) string {
	if v := viper.GetString(key); v != "" {
		return v
	}

	return defaultVal
}

// getIntOrDefault returns the viper int value for key, or defaultVal when zero or unset.
func getIntOrDefault(key string, defaultVal int) int {
	if v := viper.GetInt(key); v > 0 {
		return v
	}

	return defaultVal
}

// parseCORSOrigins reads a comma-separated list of allowed CORS origins.
// Falls back to ["*"] when the key is unset.
func parseCORSOrigins(key string) []string {
	raw := viper.GetString(key)
	if raw == "" {
		return []string{appconst.DefaultAllowedOrigins}
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			origins = append(origins, o)
		}
	}

	if len(origins) == 0 {
		return []string{appconst.DefaultAllowedOrigins}
	}

	return origins
}
