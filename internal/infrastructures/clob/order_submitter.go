package clob

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	"github.com/darmayasa221/polymarket-go/internal/domains/order"
)

// Compile-time assertion: OrderSubmitter implements tradingports.OrderSubmitter.
var _ tradingports.OrderSubmitter = (*OrderSubmitter)(nil)

// OrderSubmitter submits signed orders to the Polymarket CLOB.
type OrderSubmitter struct{ client *Client }

// NewOrderSubmitter creates an OrderSubmitter.
func NewOrderSubmitter(client *Client) *OrderSubmitter { return &OrderSubmitter{client: client} }

type submitOrderRequest struct {
	MarketID      string `json:"marketId"`
	TokenID       string `json:"tokenId"`
	Side          int    `json:"side"`
	Outcome       string `json:"outcome"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	Type          string `json:"type"`
	FeeRateBps    uint64 `json:"feeRateBps"`
	SignatureType uint8  `json:"signatureType"`
	Signature     string `json:"signature"`
}

type submitOrderResponse struct {
	OrderID string `json:"orderID"`
}

// Submit posts a signed order to the CLOB and returns the CLOB-assigned orderID.
// signature is the 65-byte EIP-712 signature from the interfaces layer.
func (s *OrderSubmitter) Submit(ctx context.Context, o *order.Order, signature []byte) (string, error) {
	body := submitOrderRequest{
		MarketID:      o.MarketID(),
		TokenID:       string(o.TokenID()),
		Side:          int(o.Side()),
		Outcome:       string(o.Outcome()),
		Price:         o.Price().String(),
		Size:          o.Size().String(),
		Type:          string(o.Type()),
		FeeRateBps:    o.FeeRateBps(),
		SignatureType: o.SignatureType(),
		Signature:     "0x" + hex.EncodeToString(signature),
	}
	var resp submitOrderResponse
	if err := s.client.do(ctx, http.MethodPost, "/order", body, &resp); err != nil {
		return "", fmt.Errorf("order submitter: submit: %w", err)
	}
	return resp.OrderID, nil
}

// Cancel sends DELETE /order/:clobOrderID to cancel an active order.
func (s *OrderSubmitter) Cancel(ctx context.Context, clobOrderID string) error {
	return s.client.do(ctx, http.MethodDelete, "/order/"+clobOrderID, nil, nil)
}
