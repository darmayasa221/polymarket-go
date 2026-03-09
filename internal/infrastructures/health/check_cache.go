package health

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/health"
	healthtypes "github.com/darmayasa221/polymarket-go/internal/applications/health/types"
	healthconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/health"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
)

// Compile-time assertion: CacheChecker implements health.Component.
var _ health.Component = (*CacheChecker)(nil)

// CacheChecker checks Redis connectivity.
type CacheChecker struct {
	cache *redisclient.Client
}

// NewCacheChecker creates a new cache health checker.
func NewCacheChecker(cache *redisclient.Client) *CacheChecker {
	return &CacheChecker{cache: cache}
}

// Name returns the component name.
func (c *CacheChecker) Name() string { return "cache" }

// Check verifies the cache is reachable.
func (c *CacheChecker) Check(ctx context.Context) healthtypes.ComponentStatus {
	if err := c.cache.Ping(ctx); err != nil {
		return healthtypes.ComponentStatus{
			Name:    c.Name(),
			Status:  healthconst.StatusUnhealthy,
			Message: err.Error(),
		}
	}
	return healthtypes.ComponentStatus{Name: c.Name(), Status: healthconst.StatusHealthy}
}
