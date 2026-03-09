package token

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// Params holds input parameters for creating a new Token.
type Params struct {
	UserID    string
	Value     string
	Type      string
	Purpose   string
	ExpiresAt time.Time
}

// New creates and validates a new Token aggregate.
// This is the ONLY way to create a Token — never use struct literal.
func New(p Params) (*Token, error) {
	t := &Token{
		id:        valueobjects.NewID(),
		userID:    valueobjects.ID(p.UserID),
		value:     TokenValue(p.Value),
		tokenType: p.Type,
		purpose:   p.Purpose,
		expiresAt: p.ExpiresAt,
		createdAt: timeutil.Now(),
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	return t, nil
}

// ReconstitutedParams holds data for rebuilding a Token from persistence.
type ReconstitutedParams struct {
	ID        string
	UserID    string
	Value     string
	Type      string
	Purpose   string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// Reconstitute rebuilds a Token from persisted data.
// Use ONLY in repository implementations when loading from DB.
func Reconstitute(p ReconstitutedParams) *Token {
	return &Token{
		id:        valueobjects.ID(p.ID),
		userID:    valueobjects.ID(p.UserID),
		value:     TokenValue(p.Value),
		tokenType: p.Type,
		purpose:   p.Purpose,
		expiresAt: p.ExpiresAt,
		createdAt: p.CreatedAt,
	}
}
