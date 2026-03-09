package logging

import (
	"errors"
	"time"
)

// statusCoder is a local duck-type interface so commons/logging never imports
// commons/errors/types. Any error with GetHTTPStatus() satisfies it — including
// DomainError — allowing log-level routing without a cross-layer import.
type statusCoder interface {
	GetHTTPStatus() int
}

// Operation logs the start and end of an operation with duration.
type Operation struct {
	logger *Logger
	name   string
	start  time.Time
}

// StartOperation begins timing an operation and logs its start.
func StartOperation(l *Logger, name string) *Operation {
	l.Info(name+" started", FieldOperation(name))
	return &Operation{logger: l, name: name, start: time.Now()}
}

// Success logs a successful operation completion with duration.
func (o *Operation) Success() {
	ms := float64(time.Since(o.start).Milliseconds())
	o.logger.Info(o.name+" completed", FieldOperation(o.name), FieldDuration(ms))
}

// Failure logs a failed operation with error and duration.
// 4xx errors (client errors) are logged at WARN; 5xx and unknown errors at ERROR.
func (o *Operation) Failure(err error) {
	ms := float64(time.Since(o.start).Milliseconds())
	var sc statusCoder
	if errors.As(err, &sc) && sc.GetHTTPStatus() < 500 {
		o.logger.Warn(o.name+" failed", FieldOperation(o.name), FieldDuration(ms), FieldError(err))
		return
	}
	o.logger.Error(o.name+" failed", FieldOperation(o.name), FieldDuration(ms), FieldError(err))
}
