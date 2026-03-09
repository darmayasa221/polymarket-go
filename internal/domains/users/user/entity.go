// Package user defines the User aggregate — the core entity of the users bounded context.
package user

import (
	"time"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// UserID is a typed value object for user identity.
// Never use raw string for user IDs.
type UserID = valueobjects.ID

// Email is a typed value object for email addresses.
type Email string

// String returns the string representation.
func (e Email) String() string { return string(e) }

// HashedPassword is a typed value object for hashed passwords.
// Never store or compare plain text passwords as this type.
type HashedPassword string

// String returns the string representation.
func (p HashedPassword) String() string { return string(p) }

// User is the User aggregate root.
// Created only via New() — never use struct literal directly.
type User struct {
	id             UserID
	username       string
	email          Email
	hashedPassword HashedPassword
	fullName       string
	createdAt      time.Time
	updatedAt      time.Time
}

// ID returns the user's identity.
func (u *User) ID() UserID { return u.id }

// Username returns the user's username.
func (u *User) Username() string { return u.username }

// Email returns the user's email.
func (u *User) Email() Email { return u.email }

// HashedPassword returns the user's hashed password.
func (u *User) HashedPassword() HashedPassword { return u.hashedPassword }

// FullName returns the user's full name.
func (u *User) FullName() string { return u.fullName }

// CreatedAt returns when the user was created.
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt returns when the user was last updated.
func (u *User) UpdatedAt() time.Time { return u.updatedAt }
