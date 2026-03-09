// Package events defines domain events for the User aggregate.
// Events are emitted when significant state changes occur.
package events

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// UserRegistered is emitted when a new user successfully registers.
type UserRegistered struct {
	UserID     valueobjects.ID
	Username   string
	Email      string
	OccurredAt time.Time
}

// NewUserRegistered creates a UserRegistered event.
func NewUserRegistered(userID valueobjects.ID, username, email string) UserRegistered {
	return UserRegistered{
		UserID:     userID,
		Username:   username,
		Email:      email,
		OccurredAt: timeutil.Now(),
	}
}

// EventName returns the event identifier.
func (e UserRegistered) EventName() string { return NameUserRegistered }
