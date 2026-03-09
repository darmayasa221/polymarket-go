// Package events defines the domain event dispatcher interface.
// The implementation lives in infrastructures/ and is wired in the container.
package events

import "context"

// Dispatcher routes domain events to their registered subscribers.
type Dispatcher interface {
	// Dispatch delivers the event to all registered subscribers.
	// Returns an error if dispatching fails for any subscriber.
	Dispatch(ctx context.Context, event any) error
}
