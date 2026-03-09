// Package constants defines HTTP-layer constants.
package constants

// Middleware error response message strings.
// These appear verbatim in the "error" field of middleware-generated responses.
const (
	MsgAuthRequired      = "authentication required"
	MsgRequestTimedOut   = "request timed out"
	MsgBodyTooLarge      = "request body too large"
	MsgInternalError     = "an unexpected error occurred"
	MsgRateLimitExceeded = "rate limit exceeded"

	// KeyPrefixRateLimit is the Redis key prefix for per-IP rate limit counters.
	KeyPrefixRateLimit = "ratelimit:"
)
