// Package clob provides the Polymarket CLOB HTTP client with L2 HMAC-SHA256 auth.
package clob

// Config holds credentials for the Polymarket CLOB API.
type Config struct {
	// BaseURL is the CLOB API base URL (e.g. "https://clob.polymarket.com").
	BaseURL string
	// Address is the Ethereum EOA wallet address (POLY_ADDRESS header).
	Address string
	// APIKey is the CLOB API key (POLY_API_KEY header).
	APIKey string
	// APISecret is the base64-encoded CLOB API secret used for HMAC signing.
	APISecret string
	// APIPassphrase is the CLOB API passphrase (POLY_API_PASSPHRASE header).
	APIPassphrase string
}
