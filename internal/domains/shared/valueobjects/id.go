// Package valueobjects defines shared domain value objects.
// Value objects have no identity — they are defined by their value.
package valueobjects

import "github.com/darmayasa221/polymarket-go/internal/commons/crypto"

// ID is a typed domain identifier. Never use raw string for entity IDs.
type ID string

// NewID generates a new random ID.
func NewID() ID {
	return ID(crypto.GenerateUUID())
}

// String returns the string representation.
func (id ID) String() string {
	return string(id)
}

// IsEmpty returns true if the ID is empty.
func (id ID) IsEmpty() bool {
	return string(id) == ""
}

// Equals returns true if both IDs have the same value.
func (id ID) Equals(other ID) bool {
	return id == other
}
