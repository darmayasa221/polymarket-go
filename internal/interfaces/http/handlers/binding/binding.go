// Package binding provides shared HTTP binding error mapping for handlers.
package binding

import (
	"errors"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"

	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// MapError converts a Gin binding/validation error into a domain error.
// If err contains go-playground/validator ValidationErrors, a ValidationError
// with per-field violations is returned. Otherwise a generic ClientError is returned.
func MapError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		violations := make([]types.FieldViolation, 0, len(ve))
		for _, fe := range ve {
			violations = append(violations, types.FieldViolation{
				Field:   toSnakeCase(fe.Field()),
				Message: tagMessage(fe.Tag()),
			})
		}
		return types.NewValidationErrorWithMessage(errkeys.ErrValidationFailed, "validation failed", violations)
	}
	return types.NewClientError(errkeys.ErrInvalidRequestBody)
}

// tagMessage converts a validator field tag into a human-readable message.
func tagMessage(tag string) string {
	switch tag {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "is too short"
	case "max":
		return "is too long"
	case "len":
		return "has an invalid length"
	case "oneof":
		return "is not one of the allowed values"
	case "uuid":
		return "must be a valid UUID"
	case "url":
		return "must be a valid URL"
	default:
		return "is invalid"
	}
}

// toSnakeCase converts a PascalCase or camelCase identifier to snake_case.
// Examples: "FullName" → "full_name", "Email" → "email", "UserID" → "user_id".
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	var b strings.Builder
	b.Grow(len(s) + 4)

	for i, r := range runes {
		if !unicode.IsUpper(r) {
			b.WriteRune(r)
			continue
		}
		// Insert underscore before this uppercase rune when:
		//   (a) it is not the first character, AND
		//   (b) the previous character was lowercase  — e.g. "Name" in "FullName"
		//       OR the next character is lowercase    — e.g. second "U" in "UserID"→ "user_id"
		insertUnderscore := i > 0 &&
			(unicode.IsLower(runes[i-1]) ||
				(i+1 < len(runes) && unicode.IsLower(runes[i+1])))
		if insertUnderscore {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
