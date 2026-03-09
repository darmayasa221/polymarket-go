// Package container provides dependency injection (Facade Pattern).
// The Container is a Facade — one entry point to access all wired dependencies.
package container

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth"
	"github.com/darmayasa221/polymarket-go/internal/applications/health"
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers"
)

// Container holds all application dependencies.
// Built once at startup — never rebuild at runtime.
type Container struct {
	// Use cases
	AddUser     adduser.UseCase
	GetUser     getuser.UseCase
	ListUsers   listusers.UseCase
	LoginUser   loginuser.UseCase
	LogoutUser  logoutuser.UseCase
	RefreshAuth refreshauth.UseCase

	// Security
	TokenManager security.TokenManager

	// Health
	HealthChecker *health.Checker
}
