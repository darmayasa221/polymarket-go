package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/position"
)

const errPositionNotFound = "PORTFOLIO.POSITION_NOT_FOUND"
const errPositionGetFailed = "PORTFOLIO.GET_FAILED"

func scanPosition(row *sql.Row) (*position.Position, error) {
	var (
		id, asset, tokenID, outcome, size, avgPrice, marketID string
		openedAt                                              sql.NullTime
		closedAt                                              sql.NullTime
	)
	if err := row.Scan(&id, &asset, &tokenID, &outcome, &size, &avgPrice, &marketID, &openedAt, &closedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NewNotFoundError(errPositionNotFound)
		}
		return nil, errtypes.NewInternalServerError(errPositionGetFailed)
	}
	return buildPosition(id, asset, tokenID, outcome, size, avgPrice, marketID, openedAt.Time, closedAt)
}

func scanPositions(rows *sql.Rows) ([]*position.Position, error) {
	var positions []*position.Position
	for rows.Next() {
		var (
			id, asset, tokenID, outcome, size, avgPrice, marketID string
			openedAt                                              sql.NullTime
			closedAt                                              sql.NullTime
		)
		if err := rows.Scan(&id, &asset, &tokenID, &outcome, &size, &avgPrice, &marketID, &openedAt, &closedAt); err != nil {
			return nil, errtypes.NewInternalServerError(errPositionGetFailed)
		}
		pos, err := buildPosition(id, asset, tokenID, outcome, size, avgPrice, marketID, openedAt.Time, closedAt)
		if err != nil {
			return nil, err
		}
		positions = append(positions, pos)
	}
	if err := rows.Err(); err != nil {
		return nil, errtypes.NewInternalServerError(errPositionGetFailed)
	}
	return positions, nil
}

func buildPosition(id, asset, tokenID, outcome, sizeStr, avgPriceStr, marketID string, openedAt time.Time, closedAt sql.NullTime) (*position.Position, error) {
	sizeDec, err := decimal.NewFromString(sizeStr)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errPositionGetFailed)
	}
	avgPriceDec, err := decimal.NewFromString(avgPriceStr)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errPositionGetFailed)
	}
	var closedAtPtr *time.Time
	if closedAt.Valid {
		closedAtPtr = &closedAt.Time
	}
	return position.Reconstitute(position.ReconstitutedParams{
		ID:       id,
		Asset:    market.Asset(asset),
		TokenID:  polyid.TokenID(tokenID),
		Outcome:  market.Outcome(outcome),
		Size:     sizeDec,
		AvgPrice: avgPriceDec,
		MarketID: marketID,
		OpenedAt: openedAt,
		ClosedAt: closedAtPtr,
	}), nil
}
