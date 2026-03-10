package postgres

import (
	"context"

	"github.com/darmayasa221/polymarket-go/internal/commons/crypto"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

const errSavePriceFailed = "PRICING.SAVE_FAILED"

// Save persists a price observation.
// A UUID surrogate key is generated because oracle.Price has no ID().
func (r *Repository) Save(ctx context.Context, price *oracle.Price) error {
	const query = `
		INSERT INTO prices (id, asset, source, value, rounded_at, received_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	var roundedAt any
	if !price.RoundedAt().IsZero() {
		roundedAt = price.RoundedAt()
	}

	_, err := r.db.ExecContext(ctx, query,
		crypto.GenerateUUID(),
		price.Asset(),
		string(price.Source()),
		price.Value().String(),
		roundedAt,
		price.ReceivedAt(),
	)
	if err != nil {
		return errtypes.NewInternalServerError(errSavePriceFailed)
	}
	return nil
}
