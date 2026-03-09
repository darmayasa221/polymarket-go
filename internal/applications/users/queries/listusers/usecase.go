package listusers

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// Compile-time assertion.
var _ UseCase = (*useCase)(nil)

type useCase struct {
	userRepo repository.User
}

// New creates a ListUsers use case.
func New(userRepo repository.User) UseCase {
	return &useCase{userRepo: userRepo}
}

// ExecuteOffset lists users using offset pagination.
func (uc *useCase) ExecuteOffset(ctx context.Context, input dto.Input) (dto.OffsetOutput, error) {
	if input.OffsetParams.Page <= 0 {
		return dto.OffsetOutput{}, types.NewClientError(ErrInvalidPage)
	}
	if input.OffsetParams.PageSize <= 0 {
		return dto.OffsetOutput{}, types.NewClientError(ErrInvalidPageSize)
	}

	result, err := uc.userRepo.ListOffset(ctx, input.OffsetParams)
	if err != nil {
		return dto.OffsetOutput{}, err
	}
	return dto.OffsetOutput{
		Users:      toUserItems(result.Items),
		TotalItems: result.TotalItems,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// ExecuteCursor lists users using cursor pagination.
func (uc *useCase) ExecuteCursor(ctx context.Context, input dto.Input) (dto.CursorOutput, error) {
	if input.CursorParams.PageSize <= 0 {
		return dto.CursorOutput{}, types.NewClientError(ErrInvalidPageSize)
	}

	result, err := uc.userRepo.ListCursor(ctx, input.CursorParams)
	if err != nil {
		return dto.CursorOutput{}, err
	}
	return dto.CursorOutput{
		Users:      toUserItems(result.Items),
		NextCursor: result.NextCursor,
		PrevCursor: result.PrevCursor,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}, nil
}

func toUserItems(users []*user.User) []dto.UserItem {
	items := make([]dto.UserItem, len(users))
	for i, u := range users {
		items[i] = dto.UserItem{
			ID:        u.ID().String(),
			Username:  u.Username(),
			Email:     u.Email().String(),
			FullName:  u.FullName(),
			CreatedAt: u.CreatedAt(),
		}
	}
	return items
}
