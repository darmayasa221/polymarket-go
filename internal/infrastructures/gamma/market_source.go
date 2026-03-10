package gamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	marketwatchports "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/ports"
	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
	"github.com/darmayasa221/polymarket-go/internal/domains/market"
)

// Compile-time assertion: MarketSource implements marketwatchports.MarketSource.
var _ marketwatchports.MarketSource = (*MarketSource)(nil)

// trackedAssets are the four tickers the bot trades.
var trackedAssets = []string{"btc", "eth", "sol", "xrp"}

// MarketSource fetches active 5-minute markets from the Polymarket Gamma API.
type MarketSource struct {
	cfg    Config
	client *http.Client
}

// NewMarketSource creates a MarketSource.
func NewMarketSource(cfg Config) *MarketSource {
	return &MarketSource{cfg: cfg, client: &http.Client{Timeout: 10 * time.Second}}
}

// gammaInnerMarket is the nested market object inside a Gamma API event response.
// Fields confirmed against https://docs.polymarket.com/api-reference/events/list-events
type gammaInnerMarket struct {
	ID                string  `json:"id"`
	ConditionID       string  `json:"conditionId"`
	Slug              string  `json:"slug"`
	Outcomes          string  `json:"outcomes"`     // JSON-encoded string: "[\"Up\",\"Down\"]"
	ClobTokenIDs      string  `json:"clobTokenIds"` // JSON-encoded string: "[\"id1\",\"id2\"]"
	OrderPriceMinTick float64 `json:"orderPriceMinTickSize"`
	EnableOrderBook   bool    `json:"enableOrderBook"`
	Active            bool    `json:"active"`
	Closed            bool    `json:"closed"`
}

// gammaEvent is the top-level event returned by GET /events?slug=...
type gammaEvent struct {
	ID      string             `json:"id"`
	Slug    string             `json:"slug"`
	Closed  bool               `json:"closed"`
	Markets []gammaInnerMarket `json:"markets"`
}

// FetchActive5mMarkets calls GET /events?slug=btc-updown-5m-{ts}&slug=eth-... for the
// current window boundary and returns open markets for all tracked assets.
func (s *MarketSource) FetchActive5mMarkets(ctx context.Context) ([]*market.Market, error) {
	slugs := currentWindowSlugs()
	apiURL := s.cfg.BaseURL + "/events?" + buildSlugQuery(slugs)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, http.NoBody)
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
		if e.Closed || len(e.Markets) == 0 {
			continue
		}
		m, err := buildMarket(e)
		if err != nil {
			continue // skip unparseable events rather than aborting all
		}
		markets = append(markets, m)
	}

	if len(markets) == 0 {
		return nil, fmt.Errorf("gamma: no active 5m markets found for slugs %v", slugs)
	}
	return markets, nil
}

// Slug pattern: {ticker}-updown-5m-{floor(unix/300)*300}.
func currentWindowSlugs() []string {
	now := timeutil.Now().Unix()
	boundary := now - (now % 300)
	slugs := make([]string, len(trackedAssets))
	for i, asset := range trackedAssets {
		slugs[i] = fmt.Sprintf("%s-updown-5m-%d", asset, boundary)
	}
	return slugs
}

// buildSlugQuery constructs query string "slug=btc-...&slug=eth-...".
func buildSlugQuery(slugs []string) string {
	vals := make(url.Values)
	for _, s := range slugs {
		vals.Add("slug", s)
	}
	return vals.Encode()
}

// buildMarket converts a gammaEvent into a domain market.Market via Reconstitute.
func buildMarket(e gammaEvent) (*market.Market, error) {
	inner := e.Markets[0]

	asset, err := assetFromSlug(inner.Slug)
	if err != nil {
		return nil, err
	}
	windowStart, err := windowStartFromSlug(inner.Slug)
	if err != nil {
		return nil, err
	}
	upTokenID, downTokenID, err := tokenIDsFromClobIDs(inner.ClobTokenIDs, inner.Outcomes)
	if err != nil {
		return nil, err
	}

	tickSize := decimal.NewFromFloat(inner.OrderPriceMinTick)
	if tickSize.IsZero() {
		tickSize = decimal.NewFromFloat(0.01) // safe default
	}

	return market.Reconstitute(market.ReconstitutedParams{
		ID:          inner.ID,
		Asset:       asset,
		WindowStart: windowStart,
		ConditionID: polyid.ConditionID(inner.ConditionID),
		UpTokenID:   upTokenID,
		DownTokenID: downTokenID,
		TickSize:    tickSize,
		FeeEnabled:  inner.EnableOrderBook,
		Active:      inner.Active && !inner.Closed,
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

// tokenIDsFromClobIDs extracts Up/Down token IDs from the JSON-encoded clobTokenIds string.
// outcomes is also JSON-encoded: "[\"Up\",\"Down\"]"
// clobTokenIds order matches outcomes order.
func tokenIDsFromClobIDs(clobTokenIDsJSON, outcomesJSON string) (up, down polyid.TokenID, err error) {
	var tokenIDs []string
	if err = json.Unmarshal([]byte(clobTokenIDsJSON), &tokenIDs); err != nil {
		return "", "", fmt.Errorf("parse clobTokenIds: %w", err)
	}
	var outcomes []string
	if err = json.Unmarshal([]byte(outcomesJSON), &outcomes); err != nil {
		return "", "", fmt.Errorf("parse outcomes: %w", err)
	}
	if len(tokenIDs) != len(outcomes) || len(tokenIDs) < 2 {
		return "", "", fmt.Errorf("mismatched tokenIds/outcomes lengths: %d vs %d", len(tokenIDs), len(outcomes))
	}
	for i, outcome := range outcomes {
		switch outcome {
		case "Up":
			up = polyid.TokenID(tokenIDs[i])
		case "Down":
			down = polyid.TokenID(tokenIDs[i])
		}
	}
	if up == "" || down == "" {
		return "", "", fmt.Errorf("missing Up or Down token in outcomes %v", outcomes)
	}
	return up, down, nil
}
