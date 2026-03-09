package usecases

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth"
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	authrepo "github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	userrepo "github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
)

// ProvideLoginUser wires the LoginUser use case.
func ProvideLoginUser(
	userRepo userrepo.User,
	authRepo authrepo.Authentication,
	enc security.Encryption,
	tm security.TokenManager,
) loginuser.UseCase {
	return loginuser.New(userRepo, authRepo, enc, tm)
}

// ProvideLogoutUser wires the LogoutUser use case.
func ProvideLogoutUser(authRepo authrepo.Authentication) logoutuser.UseCase {
	return logoutuser.New(authRepo)
}

// ProvideRefreshAuth wires the RefreshAuth use case.
func ProvideRefreshAuth(
	authRepo authrepo.Authentication,
	tm security.TokenManager,
) refreshauth.UseCase {
	return refreshauth.New(authRepo, tm)
}
