package events

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// UserDeleted is emitted when a user is deleted.
type UserDeleted struct {
	UserID     valueobjects.ID
	OccurredAt time.Time
}

// NewUserDeleted creates a UserDeleted event.
func NewUserDeleted(userID valueobjects.ID) UserDeleted {
	return UserDeleted{UserID: userID, OccurredAt: timeutil.Now()}
}

// EventName returns the event identifier.
func (e UserDeleted) EventName() string { return NameUserDeleted }
