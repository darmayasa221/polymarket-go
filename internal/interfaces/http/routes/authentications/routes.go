// Package authentications registers HTTP routes for the authentications domain.
package authentications

import (
	"github.com/gin-gonic/gin"

	authhandler "github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/authentications"
)

// Register mounts all authentication routes on the given router group.
func Register(rg *gin.RouterGroup, h *authhandler.Handler) {
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.Refresh)
}
