package rules

import "unicode"

// HasUppercase returns true if the string contains at least one uppercase letter.
func HasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// HasLowercase returns true if the string contains at least one lowercase letter.
func HasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// HasDigit returns true if the string contains at least one digit.
func HasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// IsStrongPassword returns true if the password meets strength requirements:
// min 8 chars, has uppercase, lowercase, and digit.
func IsStrongPassword(s string) bool {
	return len(s) >= 8 && HasUppercase(s) && HasLowercase(s) && HasDigit(s)
}
