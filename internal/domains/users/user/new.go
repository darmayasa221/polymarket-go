package user

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// Params holds the input parameters for creating a new User.
type Params struct {
	Username       string
	Email          string
	HashedPassword string
	FullName       string
}

// New creates and validates a new User aggregate.
// This is the ONLY way to create a User — never use struct literal.
func New(p Params) (*User, error) {
	now := timeutil.Now()
	u := &User{
		id:             valueobjects.NewID(),
		username:       p.Username,
		email:          Email(p.Email),
		hashedPassword: HashedPassword(p.HashedPassword),
		fullName:       p.FullName,
		createdAt:      now,
		updatedAt:      now,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

// ReconstitutedParams holds data for rebuilding a User from persistence.
type ReconstitutedParams struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
	FullName       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Reconstitute rebuilds a User from persisted data.
// Used by repository implementations — does NOT generate new ID or validate.
func Reconstitute(p ReconstitutedParams) *User {
	return &User{
		id:             valueobjects.ID(p.ID),
		username:       p.Username,
		email:          Email(p.Email),
		hashedPassword: HashedPassword(p.HashedPassword),
		fullName:       p.FullName,
		createdAt:      p.CreatedAt,
		updatedAt:      p.UpdatedAt,
	}
}
