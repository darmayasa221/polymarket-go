package logging

import "github.com/darmayasa221/polymarket-go/internal/commons/stringutil"

// SanitizeEmail returns a masked email for safe logging.
func SanitizeEmail(email string) string { return stringutil.MaskEmail(email) }

// SanitizeToken returns a masked token for safe logging.
func SanitizeToken(token string) string { return stringutil.MaskString(token) }
