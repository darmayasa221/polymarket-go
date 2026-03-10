package clob

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	sharedports "github.com/darmayasa221/polymarket-go/internal/applications/shared/ports"
)

// Compile-time assertion: FeeRateProvider implements sharedports.FeeRateProvider.
var _ sharedports.FeeRateProvider = (*FeeRateProvider)(nil)

// FeeRateProvider fetches live fee rates from GET /fee-rate.
type FeeRateProvider struct{ client *Client }

// NewFeeRateProvider creates a FeeRateProvider.
func NewFeeRateProvider(client *Client) *FeeRateProvider {
	return &FeeRateProvider{client: client}
}

type feeRateResponse struct {
	BaseFee uint64 `json:"base_fee"`
}

// FetchFeeRate calls GET /fee-rate?token_id=<tokenID> and returns the base fee in bps.
func (p *FeeRateProvider) FetchFeeRate(ctx context.Context, tokenID string) (uint64, error) {
	path := "/fee-rate?token_id=" + tokenID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.client.cfg.BaseURL+path, http.NoBody)
	if err != nil {
		return 0, fmt.Errorf("clob fee-rate: build request: %w", err)
	}
	if err := setL2Headers(req, p.client.cfg, ""); err != nil {
		return 0, err
	}
	resp, err := p.client.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("clob fee-rate: request: %w", err)
	}
	defer resp.Body.Close()

	var result feeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("clob fee-rate: decode: %w", err)
	}
	return result.BaseFee, nil
}
