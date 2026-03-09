package token

import "time"

const (
	DefaultAccessTokenDuration  = 15 * time.Minute
	DefaultRefreshTokenDuration = 7 * 24 * time.Hour

	// String forms used as fallback defaults in string-based config parsers.
	DefaultAccessTokenDurationStr  = "15m"
	DefaultRefreshTokenDurationStr = "168h"
)
