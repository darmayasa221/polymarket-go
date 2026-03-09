// Package events defines domain events for the Token aggregate.
package events

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// TokenCreated is emitted when a new token is created (login/refresh).
type TokenCreated struct {
	TokenID    valueobjects.ID
	UserID     valueobjects.ID
	TokenType  string
	OccurredAt time.Time
}

// NewTokenCreated creates a TokenCreated event.
func NewTokenCreated(tokenID, userID valueobjects.ID, tokenType string) TokenCreated {
	return TokenCreated{
		TokenID:    tokenID,
		UserID:     userID,
		TokenType:  tokenType,
		OccurredAt: timeutil.Now(),
	}
}

// EventName returns the event identifier.
func (e TokenCreated) EventName() string { return NameTokenCreated }
