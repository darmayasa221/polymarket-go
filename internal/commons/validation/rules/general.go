// Package rules provides reusable validation rule functions.
package rules

// IsRequired returns true if the string is non-empty after trimming.
func IsRequired(s string) bool {
	return s != ""
}

// IsMinLength returns true if the string meets the minimum length.
func IsMinLength(s string, minLen int) bool {
	return len([]rune(s)) >= minLen
}

// IsMaxLength returns true if the string does not exceed the maximum length.
func IsMaxLength(s string, maxLen int) bool {
	return len([]rune(s)) <= maxLen
}

// IsInRange returns true if the string length is within [minLen, maxLen].
func IsInRange(s string, minLen, maxLen int) bool {
	l := len([]rune(s))
	return l >= minLen && l <= maxLen
}
