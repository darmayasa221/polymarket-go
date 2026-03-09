package json

import (
	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/commons/errors"
	errkeys "github.com/darmayasa221/polymarket-go/internal/commons/errors/keys"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// FieldError represents a single field-level validation failure in an error response.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// errorBody is the standard error response envelope.
type errorBody struct {
	Success          bool         `json:"success"`
	Error            string       `json:"error"`
	Code             string       `json:"code"`
	ValidationErrors []FieldError `json:"validation_errors,omitempty"`
}

// Error sends an error response with the appropriate HTTP status from the error type.
func (p *Presenter) Error(c *gin.Context, err error) {
	status := errors.HTTPStatusOf(err)
	code := errors.CodeOf(err)
	if code == "" {
		code = errkeys.ErrInternalServer
	}

	body := errorBody{
		Success: false,
		Error:   err.Error(),
		Code:    code,
	}

	if ve, ok := errors.As[*errtypes.ValidationError](err); ok {
		fieldErrs := make([]FieldError, 0, len(ve.GetViolations()))
		for _, v := range ve.GetViolations() {
			fieldErrs = append(fieldErrs, FieldError{
				Field:   v.Field,
				Message: v.Message,
			})
		}
		body.ValidationErrors = fieldErrs
	}

	c.JSON(status, body)
}
