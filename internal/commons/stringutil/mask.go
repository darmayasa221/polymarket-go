// Package stringutil provides string manipulation utilities.
package stringutil

import "strings"

// MaskEmail masks an email for safe logging: user@domain.com → u***@domain.com.
func MaskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "***"
	}
	runes := []rune(parts[0])
	return string(runes[0]) + "***@" + parts[1]
}

// MaskString masks a string, showing only the first and last character.
func MaskString(s string) string {
	runes := []rune(s)
	if len(runes) <= 2 {
		return "***"
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-2) + string(runes[len(runes)-1])
}
