package refreshauth

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth/dto"
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	tokenconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/token"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	authrepo "github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	domaintoken "github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	authRepo     authrepo.Authentication
	tokenManager security.TokenManager
}

// New creates a RefreshAuth use case.
func New(authRepo authrepo.Authentication, tokenManager security.TokenManager) UseCase {
	return &useCase{
		authRepo:     authRepo,
		tokenManager: tokenManager,
	}
}

// Execute validates the given refresh token and returns a new token pair.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.RefreshToken == "" {
		return dto.Output{}, types.NewClientError(ErrTokenRequired)
	}

	// Verify token exists and is not expired.
	existing, err := uc.authRepo.CheckToken(ctx, domaintoken.TokenValue(input.RefreshToken))
	if err != nil {
		return dto.Output{}, types.NewAuthenticationError(ErrTokenInvalid)
	}
	if existing.IsExpired() {
		return dto.Output{}, types.NewAuthenticationError(ErrTokenInvalid)
	}

	// Generate new token pair.
	pair, err := uc.tokenManager.CreateTokenPair(ctx, existing.UserID().String())
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(security.ErrTokenCreationFailed)
	}

	// Delete old refresh token.
	if err := uc.authRepo.DeleteByValue(ctx, domaintoken.TokenValue(input.RefreshToken)); err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrDeleteOldTokenFailed)
	}

	// Persist new refresh token.
	newToken, err := domaintoken.New(domaintoken.Params{
		UserID:    existing.UserID().String(),
		Value:     pair.RefreshToken,
		Type:      tokenconst.TypeRefresh,
		Purpose:   tokenconst.PurposeAuthentication,
		ExpiresAt: pair.RefreshTokenExpiresAt,
	})
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrTokenEntityFailed)
	}

	if err := uc.authRepo.Add(ctx, newToken); err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrPersistTokenFailed)
	}

	return dto.Output{
		AccessToken:           pair.AccessToken,
		RefreshToken:          pair.RefreshToken,
		AccessTokenExpiresAt:  pair.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: pair.RefreshTokenExpiresAt,
	}, nil
}
