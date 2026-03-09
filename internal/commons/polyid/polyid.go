// Package polyid provides typed identifiers for Polymarket entities.
// Using distinct types prevents accidentally passing a TokenID where a ConditionID is expected.
package polyid

// ConditionID is the market-level identifier from Polymarket's Conditional Token Framework.
// Format: 0x hex string, e.g. "0x4a1e3d8...".
type ConditionID string

// String returns the string representation.
func (c ConditionID) String() string { return string(c) }

// IsEmpty returns true if the ID is empty.
func (c ConditionID) IsEmpty() bool { return string(c) == "" }

// TokenID is the outcome-level identifier — a 256-bit decimal string.
// Each market has exactly two: one for "Up", one for "Down".
type TokenID string

// String returns the string representation.
func (t TokenID) String() string { return string(t) }

// IsEmpty returns true if the ID is empty.
func (t TokenID) IsEmpty() bool { return string(t) == "" }

// OrderID is the CLOB-level order identifier.
// Format: 0x hex string.
type OrderID string

// String returns the string representation.
func (o OrderID) String() string { return string(o) }

// IsEmpty returns true if the ID is empty.
func (o OrderID) IsEmpty() bool { return string(o) == "" }

// SlugID is the predictable market slug built from asset + timestamp.
// Format: "{ticker}-updown-5m-{floor(unix/300)*300}", e.g. "btc-updown-5m-1700000100".
type SlugID string

// String returns the string representation.
func (s SlugID) String() string { return string(s) }

// IsEmpty returns true if the slug is empty.
func (s SlugID) IsEmpty() bool { return string(s) == "" }
