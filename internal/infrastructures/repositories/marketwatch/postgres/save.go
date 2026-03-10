package postgres

import (
	"context"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

const errSaveMarketFailed = "MARKETWATCH.SAVE_FAILED"

// Save upserts a market — safe to call repeatedly for the same (asset, window_start).
func (r *Repository) Save(ctx context.Context, m *market.Market) error {
	const query = `
		INSERT INTO markets (id, slug, asset, window_start, condition_id, up_token_id, down_token_id, tick_size, fee_enabled, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (asset, window_start) DO UPDATE SET
			slug         = EXCLUDED.slug,
			condition_id = EXCLUDED.condition_id,
			up_token_id  = EXCLUDED.up_token_id,
			down_token_id= EXCLUDED.down_token_id,
			tick_size    = EXCLUDED.tick_size,
			fee_enabled  = EXCLUDED.fee_enabled,
			active       = EXCLUDED.active`

	_, err := r.db.ExecContext(ctx, query,
		m.ID(),
		string(m.Slug()),
		string(m.Asset()),
		m.WindowStart(),
		string(m.ConditionID()),
		string(m.UpTokenID()),
		string(m.DownTokenID()),
		m.TickSize().String(),
		m.FeeEnabled(),
		m.Active(),
	)
	if err != nil {
		return errtypes.NewInternalServerError(errSaveMarketFailed)
	}
	return nil
}
