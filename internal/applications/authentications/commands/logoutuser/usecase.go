package logoutuser

import (
	"context"
	"errors"

	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	authrepo "github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	domaintoken "github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	authRepo authrepo.Authentication
}

// New creates a LogoutUser use case.
func New(authRepo authrepo.Authentication) UseCase {
	return &useCase{authRepo: authRepo}
}

// Execute deletes the given token, ending the user session.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.TokenValue == "" {
		return dto.Output{}, types.NewClientError(ErrTokenRequired)
	}

	if err := uc.authRepo.DeleteByValue(ctx, domaintoken.TokenValue(input.TokenValue)); err != nil {
		var nfe *types.NotFoundError
		if errors.As(err, &nfe) {
			// Token already gone — treat as successful logout (idempotent).
			return dto.Output{}, nil
		}
		return dto.Output{}, types.NewInternalServerError(ErrDeleteFailed)
	}

	return dto.Output{}, nil
}
