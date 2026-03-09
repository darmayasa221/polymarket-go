package getuser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	userRepo repository.User
}

// New creates a GetUser use case.
func New(userRepo repository.User) UseCase {
	return &useCase{userRepo: userRepo}
}

// Execute retrieves a user by ID.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.UserID == "" {
		return dto.Output{}, types.NewClientError(ErrUserIDRequired)
	}

	u, err := uc.userRepo.GetByID(ctx, user.UserID(input.UserID))
	if err != nil {
		return dto.Output{}, err
	}

	return dto.Output{
		ID:        u.ID().String(),
		Username:  u.Username(),
		Email:     u.Email().String(),
		FullName:  u.FullName(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}, nil
}
