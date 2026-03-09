// Package timeout cancels requests that exceed the configured duration.
package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates a middleware that cancels the request context after d.
func New(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		if ctx.Err() == context.DeadlineExceeded && !c.Writer.Written() {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"success": false,
				"code":    errkeys.ErrTimeout,
				"error":   httpconst.MsgRequestTimedOut,
			})
		}
	}
}
