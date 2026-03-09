// Package repository defines the authentication repository contract.
// The interface lives in domain — implementations live in infrastructure.
package repository

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// Authentication defines all persistence operations for the Token aggregate.
// All methods accept context.Context for cancellation and tracing.
type Authentication interface {
	// Add persists a new token.
	Add(ctx context.Context, t *token.Token) error

	// CheckToken verifies a token value exists and is valid.
	CheckToken(ctx context.Context, value token.TokenValue) (*token.Token, error)

	// DeleteByUserID removes all tokens for a user (logout all).
	DeleteByUserID(ctx context.Context, userID valueobjects.ID) error

	// DeleteByValue removes a specific token by its value.
	DeleteByValue(ctx context.Context, value token.TokenValue) error
}
