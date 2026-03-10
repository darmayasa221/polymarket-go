package postgres

import (
	"database/sql"
	"errors"

	"github.com/shopspring/decimal"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

const errMarketNotFound = "MARKETWATCH.MARKET_NOT_FOUND"
const errMarketGetFailed = "MARKETWATCH.GET_FAILED"

func scanMarket(row *sql.Row) (*market.Market, error) {
	var (
		id, slug, asset, conditionID, upTokenID, downTokenID, tickSizeStr string
		feeEnabled, active                                                bool
		windowStart                                                       sql.NullTime
	)
	if err := row.Scan(&id, &slug, &asset, &windowStart, &conditionID, &upTokenID, &downTokenID, &tickSizeStr, &feeEnabled, &active); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NewNotFoundError(errMarketNotFound)
		}
		return nil, errtypes.NewInternalServerError(errMarketGetFailed)
	}
	_ = slug
	tickSize, err := decimal.NewFromString(tickSizeStr)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errMarketGetFailed)
	}
	return market.Reconstitute(market.ReconstitutedParams{
		ID:          id,
		Asset:       market.Asset(asset),
		WindowStart: windowStart.Time,
		ConditionID: polyid.ConditionID(conditionID),
		UpTokenID:   polyid.TokenID(upTokenID),
		DownTokenID: polyid.TokenID(downTokenID),
		TickSize:    tickSize,
		FeeEnabled:  feeEnabled,
		Active:      active,
	}), nil
}

func scanMarkets(rows *sql.Rows) ([]*market.Market, error) {
	var markets []*market.Market
	for rows.Next() {
		var (
			id, slug, asset, conditionID, upTokenID, downTokenID, tickSizeStr string
			feeEnabled, active                                                bool
			windowStart                                                       sql.NullTime
		)
		if err := rows.Scan(&id, &slug, &asset, &windowStart, &conditionID, &upTokenID, &downTokenID, &tickSizeStr, &feeEnabled, &active); err != nil {
			return nil, errtypes.NewInternalServerError(errMarketGetFailed)
		}
		_ = slug
		tickSize, err := decimal.NewFromString(tickSizeStr)
		if err != nil {
			return nil, errtypes.NewInternalServerError(errMarketGetFailed)
		}
		markets = append(markets, market.Reconstitute(market.ReconstitutedParams{
			ID:          id,
			Asset:       market.Asset(asset),
			WindowStart: windowStart.Time,
			ConditionID: polyid.ConditionID(conditionID),
			UpTokenID:   polyid.TokenID(upTokenID),
			DownTokenID: polyid.TokenID(downTokenID),
			TickSize:    tickSize,
			FeeEnabled:  feeEnabled,
			Active:      active,
		}))
	}
	if err := rows.Err(); err != nil {
		return nil, errtypes.NewInternalServerError(errMarketGetFailed)
	}
	return markets, nil
}
