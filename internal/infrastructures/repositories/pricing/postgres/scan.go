package postgres

import (
	"database/sql"
	"errors"

	"github.com/shopspring/decimal"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/oracle"
)

const errPriceNotFound = "PRICING.PRICE_NOT_FOUND"
const errPriceGetFailed = "PRICING.GET_FAILED"

func scanPrice(row *sql.Row) (*oracle.Price, error) {
	var (
		asset, source, value string
		roundedAt            sql.NullTime
		receivedAt           sql.NullTime
	)
	if err := row.Scan(&asset, &source, &value, &roundedAt, &receivedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NewNotFoundError(errPriceNotFound)
		}
		return nil, errtypes.NewInternalServerError(errPriceGetFailed)
	}
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errPriceGetFailed)
	}
	p, err := oracle.New(oracle.Params{
		Asset:      asset,
		Source:     oracle.PriceSource(source),
		Value:      dec,
		RoundedAt:  roundedAt.Time,
		ReceivedAt: receivedAt.Time,
	})
	if err != nil {
		return nil, errtypes.NewInternalServerError(errPriceGetFailed)
	}
	return p, nil
}
