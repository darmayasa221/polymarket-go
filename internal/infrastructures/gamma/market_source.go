package gamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	marketwatchports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Compile-time assertion: MarketSource implements marketwatchports.MarketSource.
var _ marketwatchports.MarketSource = (*MarketSource)(nil)

// tag5m is the Gamma API tag ID for 5-minute markets.
const tag5m = "102892"

// MarketSource fetches active 5-minute markets from the Polymarket Gamma API.
type MarketSource struct {
	cfg    Config
	client *http.Client
}

// NewMarketSource creates a MarketSource.
func NewMarketSource(cfg Config) *MarketSource {
	return &MarketSource{cfg: cfg, client: &http.Client{Timeout: 10 * time.Second}}
}

type gammaToken struct {
	Outcome string `json:"outcome"`
	TokenID string `json:"token_id"`
}

type gammaEvent struct {
	ID          string       `json:"id"`
	Slug        string       `json:"slug"`
	ConditionID string       `json:"conditionId"`
	Tokens      []gammaToken `json:"tokens"`
	TickSize    string       `json:"minimum_tick_size"`
	FeeEnabled  bool         `json:"fees_enabled"`
	Closed      bool         `json:"closed"`
}

// FetchActive5mMarkets calls GET /events?tag=102892&closed=false and returns open markets.
func (s *MarketSource) FetchActive5mMarkets(ctx context.Context) ([]*market.Market, error) {
	url := s.cfg.BaseURL + "/events?tag=" + tag5m + "&closed=false"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("gamma: build request: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gamma: request: %w", err)
	}
	defer resp.Body.Close()

	var events []gammaEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("gamma: decode: %w", err)
	}

	markets := make([]*market.Market, 0, len(events))
	for _, e := range events {
		if e.Closed {
			continue
		}
		m, err := buildMarket(e)
		if err != nil {
			return nil, fmt.Errorf("gamma: build market %q: %w", e.ID, err)
		}
		markets = append(markets, m)
	}
	return markets, nil
}

// buildMarket converts a gammaEvent into a domain market.Market via Reconstitute.
func buildMarket(e gammaEvent) (*market.Market, error) {
	asset, err := assetFromSlug(e.Slug)
	if err != nil {
		return nil, err
	}
	windowStart, err := windowStartFromSlug(e.Slug)
	if err != nil {
		return nil, err
	}
	upTokenID, downTokenID, err := tokenIDs(e.Tokens)
	if err != nil {
		return nil, err
	}

	tickSize, err := decimal.NewFromString(e.TickSize)
	if err != nil {
		return nil, fmt.Errorf("gamma: invalid tick_size %q: %w", e.TickSize, err)
	}

	return market.Reconstitute(market.ReconstitutedParams{
		ID:          e.ID,
		Asset:       asset,
		WindowStart: windowStart,
		ConditionID: polyid.ConditionID(e.ConditionID),
		UpTokenID:   upTokenID,
		DownTokenID: downTokenID,
		TickSize:    tickSize,
		FeeEnabled:  e.FeeEnabled,
		Active:      true,
	}), nil
}

// assetFromSlug extracts the ticker prefix from a slug like "btc-updown-5m-1741996800".
func assetFromSlug(slug string) (market.Asset, error) {
	ticker := strings.SplitN(slug, "-", 2)[0]
	switch market.Asset(ticker) {
	case market.BTC, market.ETH, market.SOL, market.XRP:
		return market.Asset(ticker), nil
	default:
		return "", fmt.Errorf("unknown asset ticker %q in slug %q", ticker, slug)
	}
}

// windowStartFromSlug parses the Unix timestamp suffix from the slug.
func windowStartFromSlug(slug string) (time.Time, error) {
	parts := strings.Split(slug, "-")
	if len(parts) < 4 { //nolint:mnd // slug has exactly 4 parts
		return time.Time{}, fmt.Errorf("unexpected slug format: %q", slug)
	}
	var ts int64
	if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &ts); err != nil {
		return time.Time{}, fmt.Errorf("parse window start from slug %q: %w", slug, err)
	}
	return time.Unix(ts, 0).UTC(), nil
}

// tokenIDs extracts Up/Down token IDs from the tokens array.
func tokenIDs(tokens []gammaToken) (up, down polyid.TokenID, err error) {
	for _, t := range tokens {
		switch t.Outcome {
		case "Up":
			up = polyid.TokenID(t.TokenID)
		case "Down":
			down = polyid.TokenID(t.TokenID)
		}
	}
	if up == "" || down == "" {
		return "", "", fmt.Errorf("missing Up or Down token in event tokens")
	}
	return up, down, nil
}
