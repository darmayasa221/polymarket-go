package listusers

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
)

// Compile-time assertion.
var _ UseCase = (*Middleware)(nil)

const (
	ActionExecuteOffset = "listusers.ExecuteOffset"
	ActionExecuteCursor = "listusers.ExecuteCursor"
)

// Middleware wraps UseCase with logging (Decorator Pattern).
type Middleware struct {
	next   UseCase
	logger *logging.Logger
}

// NewMiddleware creates a logging decorator for the ListUsers use case.
func NewMiddleware(next UseCase, logger *logging.Logger) UseCase {
	return &Middleware{next: next, logger: logger}
}

// ExecuteOffset delegates to the next use case with logging.
func (m *Middleware) ExecuteOffset(ctx context.Context, input dto.Input) (dto.OffsetOutput, error) {
	log := logging.FromContext(ctx, m.logger)
	op := logging.StartOperation(log, ActionExecuteOffset)

	out, err := m.next.ExecuteOffset(ctx, input)
	if err != nil {
		op.Failure(err)
		return dto.OffsetOutput{}, err
	}

	op.Success()
	return out, nil
}

// ExecuteCursor delegates to the next use case with logging.
func (m *Middleware) ExecuteCursor(ctx context.Context, input dto.Input) (dto.CursorOutput, error) {
	log := logging.FromContext(ctx, m.logger)
	op := logging.StartOperation(log, ActionExecuteCursor)

	out, err := m.next.ExecuteCursor(ctx, input)
	if err != nil {
		op.Failure(err)
		return dto.CursorOutput{}, err
	}

	op.Success()
	return out, nil
}
