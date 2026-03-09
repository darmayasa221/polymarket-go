package types_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

func TestNotFoundError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		code           string
		wantCode       string
		wantHTTPStatus int
		wantErrorMsg   string
	}{
		{
			name:           "code only — Error returns code",
			code:           "USER.NOT_FOUND",
			wantCode:       "USER.NOT_FOUND",
			wantHTTPStatus: http.StatusNotFound,
			wantErrorMsg:   "USER.NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewNotFoundError(tt.code)

			// Assert
			assert.Equal(t, tt.wantCode, err.GetCode())
			assert.Equal(t, tt.wantHTTPStatus, err.GetHTTPStatus())
			assert.Equal(t, tt.wantErrorMsg, err.Error())
		})
	}
}

func TestNotFoundErrorWithMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		code         string
		message      string
		wantErrorMsg string
	}{
		{
			name:         "message set — Error returns message",
			code:         "USER.NOT_FOUND",
			message:      "user not found",
			wantErrorMsg: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewNotFoundErrorWithMessage(tt.code, tt.message)

			// Assert
			assert.Equal(t, tt.code, err.GetCode())
			assert.Equal(t, http.StatusNotFound, err.GetHTTPStatus())
			assert.Equal(t, tt.wantErrorMsg, err.Error())
		})
	}
}

func TestInvariantError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		code           string
		wantHTTPStatus int
	}{
		{
			name:           "invariant error has 422 status",
			code:           "USER.INVALID_EMAIL",
			wantHTTPStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewInvariantError(tt.code)

			// Assert
			assert.Equal(t, tt.wantHTTPStatus, err.GetHTTPStatus())
			assert.Equal(t, tt.code, err.GetCode())
		})
	}
}

func TestAuthenticationError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		code           string
		wantHTTPStatus int
	}{
		{
			name:           "authentication error has 401 status",
			code:           "TOKEN.EXPIRED",
			wantHTTPStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewAuthenticationError(tt.code)

			// Assert
			assert.Equal(t, tt.wantHTTPStatus, err.GetHTTPStatus())
			assert.Equal(t, tt.code, err.GetCode())
		})
	}
}

func TestConflictError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		code           string
		wantHTTPStatus int
	}{
		{
			name:           "conflict error has 409 status",
			code:           "USER.USERNAME_TAKEN",
			wantHTTPStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewConflictError(tt.code)

			// Assert
			assert.Equal(t, tt.wantHTTPStatus, err.GetHTTPStatus())
			assert.Equal(t, tt.code, err.GetCode())
		})
	}
}

func TestValidationError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		code               string
		violations         []types.FieldViolation
		wantHTTPStatus     int
		wantViolationCount int
		wantFirstField     string
		wantFirstMessage   string
		wantErrorMsg       string
	}{
		{
			name: "validation error has 422 status with violations",
			code: "GENERAL.VALIDATION_FAILED",
			violations: []types.FieldViolation{
				{Field: "email", Message: "must be a valid email address"},
				{Field: "full_name", Message: "must not be empty"},
			},
			wantHTTPStatus:     http.StatusUnprocessableEntity,
			wantViolationCount: 2,
			wantFirstField:     "email",
			wantFirstMessage:   "must be a valid email address",
			wantErrorMsg:       "GENERAL.VALIDATION_FAILED",
		},
		{
			name:               "validation error with empty violations",
			code:               "GENERAL.VALIDATION_FAILED",
			violations:         []types.FieldViolation{},
			wantHTTPStatus:     http.StatusUnprocessableEntity,
			wantViolationCount: 0,
			wantErrorMsg:       "GENERAL.VALIDATION_FAILED",
		},
		{
			name:               "nil violations normalised to empty slice",
			code:               "GENERAL.VALIDATION_FAILED",
			violations:         nil,
			wantHTTPStatus:     http.StatusUnprocessableEntity,
			wantViolationCount: 0,
			wantErrorMsg:       "GENERAL.VALIDATION_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewValidationError(tt.code, tt.violations)

			// Assert
			assert.Equal(t, tt.wantHTTPStatus, err.GetHTTPStatus())
			assert.Equal(t, tt.code, err.GetCode())
			assert.Equal(t, tt.wantErrorMsg, err.Error())
			assert.Len(t, err.GetViolations(), tt.wantViolationCount)

			if tt.wantViolationCount > 0 {
				assert.Equal(t, tt.wantFirstField, err.GetViolations()[0].Field)
				assert.Equal(t, tt.wantFirstMessage, err.GetViolations()[0].Message)
			}
		})
	}
}

func TestValidationErrorWithMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		code               string
		message            string
		violations         []types.FieldViolation
		wantHTTPStatus     int
		wantViolationCount int
		wantErrorMsg       string
	}{
		{
			name:    "message set — Error returns message",
			code:    "GENERAL.VALIDATION_FAILED",
			message: "validation failed",
			violations: []types.FieldViolation{
				{Field: "email", Message: "must be a valid email address"},
			},
			wantHTTPStatus:     http.StatusUnprocessableEntity,
			wantViolationCount: 1,
			wantErrorMsg:       "validation failed",
		},
		{
			name:               "nil violations normalised to empty slice",
			code:               "GENERAL.VALIDATION_FAILED",
			message:            "validation failed",
			violations:         nil,
			wantHTTPStatus:     http.StatusUnprocessableEntity,
			wantViolationCount: 0,
			wantErrorMsg:       "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange + Act
			err := types.NewValidationErrorWithMessage(tt.code, tt.message, tt.violations)

			// Assert
			assert.Equal(t, tt.wantHTTPStatus, err.GetHTTPStatus())
			assert.Equal(t, tt.code, err.GetCode())
			assert.Equal(t, tt.wantErrorMsg, err.Error())
			assert.Len(t, err.GetViolations(), tt.wantViolationCount)
		})
	}
}
