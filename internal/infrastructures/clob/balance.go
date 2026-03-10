package clob

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	sharedports "github.com/darmayasa221/polymarket-go/internal/applications/shared/ports"
)

// Compile-time assertion: BalanceProvider implements sharedports.BalanceProvider.
var _ sharedports.BalanceProvider = (*BalanceProvider)(nil)

// BalanceProvider fetches the USDC.e collateral balance from GET /balance-allowance.
type BalanceProvider struct{ client *Client }

// NewBalanceProvider creates a BalanceProvider.
func NewBalanceProvider(client *Client) *BalanceProvider {
	return &BalanceProvider{client: client}
}

type balanceResponse struct {
	Balance string `json:"balance"`
}

// FetchBalance calls GET /balance-allowance?asset_type=COLLATERAL
// and returns the USDC.e collateral balance in dollars (e.g. 10.50).
func (b *BalanceProvider) FetchBalance(ctx context.Context) (decimal.Decimal, error) {
	path := "/balance-allowance?asset_type=COLLATERAL"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.client.cfg.BaseURL+path, http.NoBody)
	if err != nil {
		return decimal.Zero, fmt.Errorf("clob balance: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "@polymarket/clob-client")
	req.Header.Set("Accept", "*/*")
	if err := setL2Headers(req, b.client.cfg, ""); err != nil {
		return decimal.Zero, err
	}

	resp, err := b.client.http.Do(req)
	if err != nil {
		return decimal.Zero, fmt.Errorf("clob balance: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("clob balance: unexpected status %d", resp.StatusCode)
	}

	var result balanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return decimal.Zero, fmt.Errorf("clob balance: decode: %w", err)
	}

	bal, err := decimal.NewFromString(result.Balance)
	if err != nil {
		return decimal.Zero, fmt.Errorf("clob balance: parse %q: %w", result.Balance, err)
	}
	return bal, nil
}
