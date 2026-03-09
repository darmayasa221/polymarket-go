// Package crypto provides cryptographic utility functions.
package crypto

import "github.com/google/uuid"

// GenerateUUID generates a new random UUID string.
// Always use this for entity IDs — never use raw uuid.New() directly.
func GenerateUUID() string {
	return uuid.New().String()
}
