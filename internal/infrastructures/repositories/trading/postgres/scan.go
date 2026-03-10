package postgres

import (
	"database/sql"
	"errors"

	"github.com/shopspring/decimal"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

const errOrderNotFound = "TRADING.ORDER_NOT_FOUND"
const errOrderGetFailed = "TRADING.GET_FAILED"

func scanOrder(row *sql.Row) (*order.Order, error) {
	var (
		id, marketID, tokenID, outcome, price, size, orderType, status string
		side, signatureType                                            int
		feeRateBps                                                     int64
		expiration                                                     sql.NullTime
		createdAt                                                      sql.NullTime
	)
	if err := row.Scan(&id, &marketID, &tokenID, &side, &outcome, &price, &size, &orderType, &expiration, &feeRateBps, &signatureType, &status, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NewNotFoundError(errOrderNotFound)
		}
		return nil, errtypes.NewInternalServerError(errOrderGetFailed)
	}
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errOrderGetFailed)
	}
	sizeDec, err := decimal.NewFromString(size)
	if err != nil {
		return nil, errtypes.NewInternalServerError(errOrderGetFailed)
	}
	return order.Reconstitute(order.ReconstitutedParams{
		ID:            polyid.OrderID(id),
		MarketID:      marketID,
		TokenID:       polyid.TokenID(tokenID),
		Side:          order.Side(side), //nolint:gosec // db value is 0 or 1
		Outcome:       market.Outcome(outcome),
		Price:         priceDec,
		Size:          sizeDec,
		Type:          order.OrderType(orderType),
		Expiration:    expiration.Time,
		FeeRateBps:    uint64(feeRateBps),   //nolint:gosec // db value is always non-negative
		SignatureType: uint8(signatureType), //nolint:gosec // db value is 0 or 1
		Status:        order.OrderStatus(status),
		CreatedAt:     createdAt.Time,
	}), nil
}

func scanOrders(rows *sql.Rows) ([]*order.Order, error) {
	var orders []*order.Order
	for rows.Next() {
		var (
			id, marketID, tokenID, outcome, price, size, orderType, status string
			side, signatureType                                            int
			feeRateBps                                                     int64
			expiration                                                     sql.NullTime
			createdAt                                                      sql.NullTime
		)
		if err := rows.Scan(&id, &marketID, &tokenID, &side, &outcome, &price, &size, &orderType, &expiration, &feeRateBps, &signatureType, &status, &createdAt); err != nil {
			return nil, errtypes.NewInternalServerError(errOrderGetFailed)
		}
		priceDec, err := decimal.NewFromString(price)
		if err != nil {
			return nil, errtypes.NewInternalServerError(errOrderGetFailed)
		}
		sizeDec, err := decimal.NewFromString(size)
		if err != nil {
			return nil, errtypes.NewInternalServerError(errOrderGetFailed)
		}
		orders = append(orders, order.Reconstitute(order.ReconstitutedParams{
			ID:            polyid.OrderID(id),
			MarketID:      marketID,
			TokenID:       polyid.TokenID(tokenID),
			Side:          order.Side(side), //nolint:gosec // db value is 0 or 1
			Outcome:       market.Outcome(outcome),
			Price:         priceDec,
			Size:          sizeDec,
			Type:          order.OrderType(orderType),
			Expiration:    expiration.Time,
			FeeRateBps:    uint64(feeRateBps),   //nolint:gosec // db value is always non-negative
			SignatureType: uint8(signatureType), //nolint:gosec // db value is 0 or 1
			Status:        order.OrderStatus(status),
			CreatedAt:     createdAt.Time,
		}))
	}
	if err := rows.Err(); err != nil {
		return nil, errtypes.NewInternalServerError(errOrderGetFailed)
	}
	return orders, nil
}
