package binding_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/binding"
)

// --- helpers to produce validator.ValidationErrors ---

type testStruct struct {
	Email    string `validate:"required,email"`
	FullName string `validate:"required"`
	UserID   string `validate:"required,uuid"`
}

func validationErrorsFor(t *testing.T, s any) validator.ValidationErrors {
	t.Helper()
	v := validator.New()
	err := v.Struct(s)
	require.Error(t, err)
	var ve validator.ValidationErrors
	require.True(t, errors.As(err, &ve))
	return ve
}

// --- MapError tests ---

func TestMapError_NonValidatorError_ReturnsClientError(t *testing.T) {
	t.Parallel()

	result := binding.MapError(errors.New("some random error"))

	var ce *types.ClientError
	require.True(t, errors.As(result, &ce), "expected *types.ClientError")
	assert.Equal(t, errkeys.ErrInvalidRequestBody, ce.GetCode())
}

func TestMapError_ValidatorError_RequiredTag(t *testing.T) {
	t.Parallel()

	ve := validationErrorsFor(t, &testStruct{})
	result := binding.MapError(ve)

	var valErr *types.ValidationError
	require.True(t, errors.As(result, &valErr), "expected *types.ValidationError")
	assert.Equal(t, errkeys.ErrValidationFailed, valErr.GetCode())
	assert.Equal(t, "validation failed", valErr.Error())

	violations := valErr.GetViolations()
	require.NotEmpty(t, violations)

	// find the "required" violation for email (first failing field)
	found := false
	for _, v := range violations {
		if v.Field == "email" && v.Message == "is required" {
			found = true
			break
		}
		if v.Field == "email" && v.Message == "must be a valid email address" {
			// email tag fires instead — still valid
			found = true
			break
		}
	}
	assert.True(t, found, "expected violation for 'email' field, got: %v", violations)
}

func TestMapError_ValidatorError_FieldsAreSnakeCase(t *testing.T) {
	t.Parallel()

	ve := validationErrorsFor(t, &testStruct{})
	result := binding.MapError(ve)

	var valErr *types.ValidationError
	require.True(t, errors.As(result, &valErr))

	for _, v := range valErr.GetViolations() {
		// field names must be lowercase (snake_case has no uppercase)
		assert.Equal(t, v.Field, lowercaseOf(v.Field),
			"field %q should be snake_case", v.Field)
	}
}

func lowercaseOf(s string) string {
	result := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// --- toSnakeCase indirect tests via MapError field output ---

func TestMapError_SnakeCaseConversions(t *testing.T) {
	t.Parallel()

	type snakeInput struct {
		FullName string `validate:"required"`
	}
	v := validator.New()
	err := v.Struct(&snakeInput{})
	require.Error(t, err)

	result := binding.MapError(err)

	var valErr *types.ValidationError
	require.True(t, errors.As(result, &valErr))

	violations := valErr.GetViolations()
	require.Len(t, violations, 1)
	assert.Equal(t, "full_name", violations[0].Field)
	assert.Equal(t, "is required", violations[0].Message)
}

func TestMapError_SnakeCaseConversions_UserID(t *testing.T) {
	t.Parallel()

	type userIDInput struct {
		UserID string `validate:"required"`
	}
	v := validator.New()
	err := v.Struct(&userIDInput{})
	require.Error(t, err)

	result := binding.MapError(err)

	var valErr *types.ValidationError
	require.True(t, errors.As(result, &valErr))

	violations := valErr.GetViolations()
	require.Len(t, violations, 1)
	assert.Equal(t, "user_id", violations[0].Field, "UserID must convert to exactly user_id (not userid)")
}
