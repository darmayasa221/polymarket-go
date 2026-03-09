package token

import "github.com/darmayasa221/polymarket-go/internal/commons/timeutil"

// IsExpired returns true if the token has expired.
func (t *Token) IsExpired() bool {
	return timeutil.Now().After(t.expiresAt)
}
