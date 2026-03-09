// Package types provides domain error types with HTTP status mapping.
package types

// DomainError is the interface that all domain error types implement.
type DomainError interface {
	error
	GetCode() string
	GetHTTPStatus() int
}

// baseError holds the common fields for all domain errors.
type baseError struct {
	code       string
	message    string
	httpStatus int
}

// Error returns the message if set, otherwise returns the code.
func (e *baseError) Error() string {
	if e.message != "" {
		return e.message
	}

	return e.code
}

// GetCode returns the domain error code.
func (e *baseError) GetCode() string {
	return e.code
}

// GetHTTPStatus returns the HTTP status code associated with this error.
func (e *baseError) GetHTTPStatus() int {
	return e.httpStatus
}
