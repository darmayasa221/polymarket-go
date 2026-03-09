// Package users registers HTTP routes for the users domain.
package users

import (
	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	usershandler "github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/users"
	authmw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/auth"
)

// Register mounts all user routes on the given router group.
func Register(rg *gin.RouterGroup, h *usershandler.Handler, tokenManager security.TokenManager) {
	users := rg.Group("/users")
	// Public
	users.POST("", h.Register)
	// Protected
	protected := users.Group("")
	protected.Use(authmw.New(tokenManager))
	protected.GET("/me", h.GetMe)
	protected.GET("/:id", h.GetByID)
	protected.GET("", h.List)
}
