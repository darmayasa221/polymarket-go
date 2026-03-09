// Package ratelimit provides IP-based rate limiting middleware using Redis.
package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"

	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// New creates an IP-based rate limiting middleware.
// maxRequests allowed per window duration.
func New(client *goredis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("%s%s", httpconst.KeyPrefixRateLimit, c.ClientIP())
		ctx := c.Request.Context()

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			if expireErr := client.Expire(ctx, key, window).Err(); expireErr != nil {
				// TTL set failed; proceed without expiry to avoid blocking the request.
				c.Next()
				return
			}
		}
		if count > int64(maxRequests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    errkeys.ErrRateLimitExceeded,
				"error":   httpconst.MsgRateLimitExceeded,
			})
			return
		}
		c.Next()
	}
}
