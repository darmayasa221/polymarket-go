package events

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// TokenRevoked is emitted when a token is revoked (logout).
type TokenRevoked struct {
	TokenID    valueobjects.ID
	UserID     valueobjects.ID
	OccurredAt time.Time
}

// NewTokenRevoked creates a TokenRevoked event.
func NewTokenRevoked(tokenID, userID valueobjects.ID) TokenRevoked {
	return TokenRevoked{
		TokenID:    tokenID,
		UserID:     userID,
		OccurredAt: timeutil.Now(),
	}
}

// EventName returns the event identifier.
func (e TokenRevoked) EventName() string { return NameTokenRevoked }
