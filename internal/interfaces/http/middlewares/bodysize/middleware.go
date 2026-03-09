// Package bodysize limits the maximum request body size.
package bodysize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates a middleware that rejects requests exceeding maxBytes.
// Falls back to httpconst.MaxRequestBodyBytes when maxBytes is zero or negative.
func New(maxBytes int64) gin.HandlerFunc {
	if maxBytes <= 0 {
		maxBytes = httpconst.MaxRequestBodyBytes
	}
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"code":    errkeys.ErrBodyTooLarge,
				"error":   httpconst.MsgBodyTooLarge,
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
