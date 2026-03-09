// Package security sets defensive HTTP security headers.
package security

import (
	"github.com/gin-gonic/gin"

	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates a middleware that sets standard security headers.
func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(httpconst.HeaderXContentTypeOptions, httpconst.SecurityNoSniff)
		c.Header(httpconst.HeaderXFrameOptions, httpconst.SecurityFrameDeny)
		c.Header(httpconst.HeaderXXSSProtection, httpconst.SecurityXSSProtection)
		c.Header(httpconst.HeaderReferrerPolicy, httpconst.SecurityReferrerPolicy)
		c.Next()
	}
}
