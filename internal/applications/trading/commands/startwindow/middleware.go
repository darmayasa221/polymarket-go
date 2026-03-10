package startwindow

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
)

// Compile-time assertion.
var _ UseCase = (*Middleware)(nil)

const actionExecute = "startwindow.Execute"

// Middleware wraps UseCase with logging (Decorator Pattern).
type Middleware struct {
	next   UseCase
	logger *logging.Logger
}

// NewMiddleware creates a logging decorator for the StartWindow use case.
func NewMiddleware(next UseCase, logger *logging.Logger) UseCase {
	return &Middleware{next: next, logger: logger}
}

// Execute delegates to the next use case with logging.
func (m *Middleware) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	log := logging.FromContext(ctx, m.logger)
	op := logging.StartOperation(log, actionExecute)

	out, err := m.next.Execute(ctx, input)
	if err != nil {
		op.Failure(err)
		return dto.Output{}, err
	}

	op.Success()
	return out, nil
}
