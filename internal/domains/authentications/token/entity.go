// Package token defines the Token aggregate for the authentications bounded context.
package token

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// TokenID is a typed value object for token identity.
type TokenID = valueobjects.ID

// TokenValue is a typed value object for the token string value.
type TokenValue string

// String returns the string representation.
func (t TokenValue) String() string { return string(t) }

// IsEmpty returns true if the token value is empty.
func (t TokenValue) IsEmpty() bool { return string(t) == "" }

// Token is the Token aggregate root.
// Created only via New() — never use struct literal.
type Token struct {
	id        TokenID
	userID    valueobjects.ID
	value     TokenValue
	tokenType string
	purpose   string
	expiresAt time.Time
	createdAt time.Time
}

// ID returns the token's identity.
func (t *Token) ID() TokenID { return t.id }

// UserID returns the associated user's ID.
func (t *Token) UserID() valueobjects.ID { return t.userID }

// Value returns the token string value.
func (t *Token) Value() TokenValue { return t.value }

// Type returns the token type (access/refresh).
func (t *Token) Type() string { return t.tokenType }

// Purpose returns the token purpose.
func (t *Token) Purpose() string { return t.purpose }

// ExpiresAt returns when the token expires.
func (t *Token) ExpiresAt() time.Time { return t.expiresAt }

// CreatedAt returns when the token was created.
func (t *Token) CreatedAt() time.Time { return t.createdAt }
