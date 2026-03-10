package postgres

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	portfolioports "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/ports"
	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

const errListClosedWithExitFailed = "PORTFOLIO.LIST_CLOSED_WITH_EXIT_FAILED"

// ListClosedWithExitPrice returns all closed positions paired with their exit price and close time.
func (r *Repository) ListClosedWithExitPrice(ctx context.Context) ([]portfolioports.ClosedPositionRecord, error) {
	const query = `SELECT id, asset, token_id, outcome, size, avg_price, market_id, opened_at, closed_at, exit_price
		FROM positions WHERE closed_at IS NOT NULL`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errListClosedWithExitFailed)
	}
	defer rows.Close()

	var records []portfolioports.ClosedPositionRecord
	for rows.Next() {
		var (
			id, asset, tokenID, outcome, sizeStr, avgPriceStr, marketID string
			openedAt, closedAt                                          sql.NullTime
			exitPriceStr                                                sql.NullString
		)
		if err := rows.Scan(&id, &asset, &tokenID, &outcome, &sizeStr, &avgPriceStr, &marketID, &openedAt, &closedAt, &exitPriceStr); err != nil {
			return nil, errtypes.NewInternalServerError(errListClosedWithExitFailed)
		}
		sizeDec, err := decimal.NewFromString(sizeStr)
		if err != nil {
			return nil, errtypes.NewInternalServerError(errListClosedWithExitFailed)
		}
		avgPriceDec, err := decimal.NewFromString(avgPriceStr)
		if err != nil {
			return nil, errtypes.NewInternalServerError(errListClosedWithExitFailed)
		}
		exitPrice, err := decimal.NewFromString(exitPriceStr.String)
		if err != nil {
			return nil, errtypes.NewInternalServerError(errListClosedWithExitFailed)
		}
		pos := position.Reconstitute(position.ReconstitutedParams{
			ID:       id,
			Asset:    market.Asset(asset),
			TokenID:  polyid.TokenID(tokenID),
			Outcome:  market.Outcome(outcome),
			Size:     sizeDec,
			AvgPrice: avgPriceDec,
			MarketID: marketID,
			OpenedAt: openedAt.Time,
			ClosedAt: &closedAt.Time,
		})
		records = append(records, portfolioports.ClosedPositionRecord{
			Pos:       pos,
			ExitPrice: exitPrice,
			ClosedAt:  closedAt.Time,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errtypes.NewInternalServerError(errListClosedWithExitFailed)
	}
	return records, nil
}
