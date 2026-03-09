package security

import "context"

// Encryption defines password hashing and comparison operations.
// Implemented in infrastructures/security/encryption/bcrypt/.
type Encryption interface {
	// Hash hashes a plain text password.
	Hash(ctx context.Context, password string) (string, error)

	// Compare verifies a plain text password against a hash.
	// Returns nil if they match, error if not.
	Compare(ctx context.Context, hashedPassword, plainPassword string) error
}
