// Package validation provides reusable validation utilities.
package validation

import "github.com/darmayasa221/polymarket-go/internal/commons/validation/types"

// Validate runs a series of named validation checks.
// Returns a MultiError if any checks fail, nil if all pass.
func Validate(checks map[string]func() (bool, string, string)) error {
	multi := &types.MultiError{}
	for field, check := range checks {
		ok, code, message := check()
		if !ok {
			multi.Add(field, code, message)
		}
	}
	if multi.HasErrors() {
		return multi
	}
	return nil
}
