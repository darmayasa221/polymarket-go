// Package health registers HTTP routes for health checks.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/applications/health"
	healthconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/health"
)

// Register mounts health check routes on the root group (no /api/v1 prefix).
func Register(rg *gin.RouterGroup, checker *health.Checker) {
	rg.GET(healthconst.PathLiveness, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": healthconst.StatusLive})
	})
	rg.GET(healthconst.PathReadiness, func(c *gin.Context) {
		report := checker.Check(c.Request.Context())
		status := http.StatusOK
		if report.Status != healthconst.StatusHealthy {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, report)
	})
}
