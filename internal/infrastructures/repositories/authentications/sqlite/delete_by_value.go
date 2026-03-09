package sqlite

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/token"
)

// DeleteByValue removes a specific token identified by its value.
// Returns NotFoundError if the token does not exist, InternalServerError on failure.
func (r *Repository) DeleteByValue(ctx context.Context, value token.TokenValue) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tokens WHERE value = ?`, value.String())
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}
	if rows == 0 {
		return errtypes.NewNotFoundError(repository.ErrTokenNotFound)
	}
	return nil
}
