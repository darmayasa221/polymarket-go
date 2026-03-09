package rules

import "regexp"

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{7,14}$`)

// IsPhone returns true if the string is a valid international phone number.
func IsPhone(s string) bool {
	return phoneRegex.MatchString(s)
}
