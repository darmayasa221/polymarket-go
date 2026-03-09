// Package health defines the health checking application service.
package health

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/health/types"
	healthconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/health"
)

// Component defines a checkable system component (DB, cache, etc).
// Implemented in infrastructures/health/.
type Component interface {
	// Check verifies the component is healthy.
	Check(ctx context.Context) types.ComponentStatus
	// Name returns the component name.
	Name() string
}

// Checker orchestrates health checks across all registered components.
type Checker struct {
	components []Component
}

// NewChecker creates a new Checker with the given components.
func NewChecker(components ...Component) *Checker {
	return &Checker{components: components}
}

// Check runs all component checks and returns the overall report.
func (c *Checker) Check(ctx context.Context) types.HealthReport {
	report := types.HealthReport{Status: healthconst.StatusHealthy}
	for _, comp := range c.components {
		status := comp.Check(ctx)
		report.Components = append(report.Components, status)
		if status.Status != healthconst.StatusHealthy {
			report.Status = healthconst.StatusUnhealthy
		}
	}
	return report
}
