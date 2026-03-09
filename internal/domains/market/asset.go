// Package market defines the Market aggregate for Polymarket 5-minute crypto markets.
package market

// Asset represents a supported cryptocurrency ticker.
type Asset string

const (
	// BTC is Bitcoin.
	BTC Asset = "btc"
	// ETH is Ethereum.
	ETH Asset = "eth"
	// SOL is Solana.
	SOL Asset = "sol"
	// XRP is Ripple.
	XRP Asset = "xrp"
)

// validAssets is the set of supported assets for fast lookup.
var validAssets = map[Asset]struct{}{
	BTC: {}, ETH: {}, SOL: {}, XRP: {},
}

// IsValid returns true if the asset is one of the four supported tickers.
func (a Asset) IsValid() bool {
	_, ok := validAssets[a]
	return ok
}
