// Package auth provides JWT authentication middleware.
package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors"
	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/telemetry"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// New creates a JWT authentication middleware.
// Validates the Bearer token and injects user_id into the Gin context.
func New(tokenManager security.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(httpconst.HeaderAuthorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, httpconst.PrefixBearer) {
			respondUnauthorized(c, errkeys.ErrUnauthorized)
			return
		}

		tokenValue := strings.TrimPrefix(authHeader, httpconst.PrefixBearer)
		claims, err := tokenManager.VerifyAccessToken(c.Request.Context(), tokenValue)
		if err != nil {
			code := errors.CodeOf(err)
			if code == "" {
				code = errkeys.ErrUnauthorized
			}
			respondUnauthorized(c, code)
			return
		}

		// Inject user ID into both Gin context (for handlers) and request context
		// (for use case middlewares that use logging.FromContext).
		ctx := telemetry.WithUserID(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(response.ContextKeyUserID, claims.UserID)
		c.Next()
	}
}

func respondUnauthorized(c *gin.Context, code string) {
	authErr := types.NewAuthenticationError(code)
	c.AbortWithStatusJSON(authErr.GetHTTPStatus(), gin.H{
		"success": false,
		"code":    code,
		"error":   httpconst.MsgAuthRequired,
	})
}
