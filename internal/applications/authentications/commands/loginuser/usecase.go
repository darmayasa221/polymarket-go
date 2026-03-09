package loginuser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	"github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	authrepo "github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	domaintoken "github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
	userrepo "github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	userRepo     userrepo.User
	authRepo     authrepo.Authentication
	encryption   security.Encryption
	tokenManager security.TokenManager
}

// New creates a LoginUser use case.
func New(
	userRepo userrepo.User,
	authRepo authrepo.Authentication,
	encryption security.Encryption,
	tokenManager security.TokenManager,
) UseCase {
	return &useCase{
		userRepo:     userRepo,
		authRepo:     authRepo,
		encryption:   encryption,
		tokenManager: tokenManager,
	}
}

// Execute authenticates a user and returns a token pair.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Username == "" {
		return dto.Output{}, types.NewClientError(ErrUsernameRequired)
	}
	if input.Password == "" {
		return dto.Output{}, types.NewClientError(ErrPasswordRequired)
	}

	// Get hashed password — do not reveal if user exists on failure.
	hashedPwd, err := uc.userRepo.GetPassword(ctx, input.Username)
	if err != nil {
		return dto.Output{}, types.NewAuthenticationError(ErrInvalidCredentials)
	}

	// Verify password.
	if err := uc.encryption.Compare(ctx, hashedPwd.String(), input.Password); err != nil {
		return dto.Output{}, types.NewAuthenticationError(ErrInvalidCredentials)
	}

	// Get user ID.
	userID, err := uc.userRepo.GetIDByUsername(ctx, input.Username)
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrGetIDFailed)
	}

	// Generate token pair.
	pair, err := uc.tokenManager.CreateTokenPair(ctx, userID.String())
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(security.ErrTokenCreationFailed)
	}

	// Persist refresh token using domain constants for type and purpose.
	refreshToken, err := domaintoken.New(domaintoken.Params{
		UserID:    userID.String(),
		Value:     pair.RefreshToken,
		Type:      token.TypeRefresh,
		Purpose:   token.PurposeAuthentication,
		ExpiresAt: pair.RefreshTokenExpiresAt,
	})
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrTokenEntityFailed)
	}

	if err := uc.authRepo.Add(ctx, refreshToken); err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrPersistTokenFailed)
	}

	return dto.Output{
		AccessToken:           pair.AccessToken,
		RefreshToken:          pair.RefreshToken,
		AccessTokenExpiresAt:  pair.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: pair.RefreshTokenExpiresAt,
		UserID:                userID.String(),
	}, nil
}
