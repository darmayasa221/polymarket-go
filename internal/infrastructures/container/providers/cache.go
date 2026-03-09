package providers

import (
	"fmt"

	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
)

// ProvideCache creates and returns a Redis client.
func ProvideCache(cfg *config.Config) (*redisclient.Client, error) {
	client, err := redisclient.New(cfg.Cache)
	if err != nil {
		return nil, fmt.Errorf("providers: cache: %w", err)
	}
	return client, nil
}
