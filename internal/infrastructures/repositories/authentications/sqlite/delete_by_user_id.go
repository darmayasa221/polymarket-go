package sqlite

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

// DeleteByUserID removes all tokens belonging to the given user.
// This operation is idempotent — deleting zero tokens is not an error.
func (r *Repository) DeleteByUserID(ctx context.Context, userID valueobjects.ID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tokens WHERE user_id = ?`, userID.String())
	if err != nil {
		return errtypes.NewInternalServerError(repository.ErrTokenDeleteFailed)
	}
	return nil
}
