package adduser

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	userRepo   repository.User
	encryption security.Encryption
}

// New creates an AddUser use case.
func New(userRepo repository.User, encryption security.Encryption) UseCase {
	return &useCase{userRepo: userRepo, encryption: encryption}
}

// Execute registers a new user.
func (uc *useCase) Execute(ctx context.Context, input dto.Input) (dto.Output, error) {
	if input.Username == "" {
		return dto.Output{}, types.NewClientError(ErrUsernameRequired)
	}
	if input.Email == "" {
		return dto.Output{}, types.NewClientError(ErrEmailRequired)
	}
	if input.Password == "" {
		return dto.Output{}, types.NewClientError(ErrPasswordRequired)
	}
	if input.FullName == "" {
		return dto.Output{}, types.NewClientError(ErrFullNameRequired)
	}

	taken, err := uc.userRepo.VerifyUsername(ctx, input.Username)
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(repository.ErrUserAddFailed)
	}
	if taken {
		return dto.Output{}, types.NewConflictError(repository.ErrUsernameTaken)
	}

	hashedPw, err := uc.encryption.Hash(ctx, input.Password)
	if err != nil {
		return dto.Output{}, types.NewInternalServerError(security.ErrPasswordHashFailed)
	}

	u, err := user.New(user.Params{
		Username:       input.Username,
		Email:          input.Email,
		HashedPassword: hashedPw,
		FullName:       input.FullName,
	})
	if err != nil {
		return dto.Output{}, err
	}

	if err := uc.userRepo.Add(ctx, u); err != nil {
		return dto.Output{}, types.NewInternalServerError(ErrAddFailed)
	}

	return dto.Output{
		ID:        u.ID().String(),
		Username:  u.Username(),
		Email:     u.Email().String(),
		FullName:  u.FullName(),
		CreatedAt: u.CreatedAt(),
	}, nil
}
