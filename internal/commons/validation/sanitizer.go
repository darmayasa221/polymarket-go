package validation

import "strings"

// SanitizeString trims spaces and normalizes whitespace.
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// SanitizeEmail lowercases and trims an email address.
func SanitizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
