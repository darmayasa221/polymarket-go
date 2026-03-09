// Package cors provides CORS middleware.
package cors

import (
	"net/http"

	gincors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates a CORS middleware using the provided allowed origins.
// Pass ["*"] for development. For production, specify exact origins (e.g. ["https://example.com"]).
func New(allowedOrigins []string) gin.HandlerFunc {
	return gincors.New(gincors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			httpconst.HeaderOrigin,
			httpconst.HeaderContentType,
			httpconst.HeaderAuthorization,
			httpconst.HeaderRequestID,
		},
		ExposeHeaders:    []string{httpconst.HeaderRequestID, httpconst.HeaderCorrelationID},
		AllowCredentials: false,
	})
}
