package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 returns the hex-encoded SHA256 hash of the input string.
func SHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
