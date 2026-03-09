// Package listusers implements the list users query use case.
package listusers

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers/dto"
)

// UseCase is the read contract for listing users with pagination.
type UseCase interface {
	// ExecuteOffset lists users using offset pagination.
	ExecuteOffset(ctx context.Context, input dto.Input) (dto.OffsetOutput, error)
	// ExecuteCursor lists users using cursor pagination.
	ExecuteCursor(ctx context.Context, input dto.Input) (dto.CursorOutput, error)
}
