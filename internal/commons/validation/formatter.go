package validation

import (
	"errors"

	"github.com/darmayasa221/polymarket-go/internal/commons/validation/types"
)

// FormatErrors extracts field errors from a MultiError into a map.
// Useful for HTTP response formatting.
func FormatErrors(err error) map[string]string {
	var multi *types.MultiError
	if !errors.As(err, &multi) {
		return nil
	}
	result := make(map[string]string, len(multi.Errors))
	for _, e := range multi.Errors {
		result[e.Field] = e.Code
	}
	return result
}
