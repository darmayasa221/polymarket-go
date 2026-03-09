package types

import "strings"

// MultiError holds multiple validation errors.
// Returned when multiple fields fail validation simultaneously.
type MultiError struct {
	Errors []*FieldError
}

// Error implements the error interface.
func (m *MultiError) Error() string {
	msgs := make([]string, len(m.Errors))
	for i, e := range m.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// Add appends a field error to the multi-error.
func (m *MultiError) Add(field, code, message string) {
	m.Errors = append(m.Errors, &FieldError{Field: field, Code: code, Message: message})
}

// HasErrors returns true if there are any validation errors.
func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}
