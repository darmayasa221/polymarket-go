// Package uow defines the Unit of Work pattern for transactional boundaries.
// Use UnitOfWork when a use case must write to multiple repositories atomically.
package uow

import "context"

// Transaction represents an active database transaction.
type Transaction interface {
	// Commit commits the transaction.
	Commit(ctx context.Context) error
	// Rollback rolls back the transaction.
	Rollback(ctx context.Context) error
}

// UnitOfWork manages transactional boundaries across repositories.
// Implemented in infrastructure layer.
type UnitOfWork interface {
	// Begin starts a new transaction and returns it.
	Begin(ctx context.Context) (Transaction, error)
}
