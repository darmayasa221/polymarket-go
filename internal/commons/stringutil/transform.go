package stringutil

import (
	"strings"
	"unicode"
)

// ToSnakeCase converts a camelCase or PascalCase string to snake_case.
// Handles acronyms correctly: "HTTPRequest" → "http_request".
func ToSnakeCase(s string) string {
	runes := []rune(s)
	var result strings.Builder
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			next := rune(0)
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			if unicode.IsLower(prev) || (next != 0 && unicode.IsLower(next)) {
				result.WriteByte('_')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// TrimSpaces trims leading and trailing spaces from a string.
func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}
