// Package recovery provides panic recovery middleware.
// Position: FIRST in the middleware chain.
package recovery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates a panic recovery middleware that logs panics and returns 500.
func New(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					logging.FieldOperation("recovery"),
					logging.FieldLayer("middleware"),
					logging.FieldError(fmt.Errorf("%v", err)),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"code":    errkeys.ErrInternalServer,
					"error":   httpconst.MsgInternalError,
				})
			}
		}()
		c.Next()
	}
}
