package dto

import "github.com/shopspring/decimal"

// Output is the result of PlaceOrder.
// The interfaces layer must sign UnsignedHash with the EOA private key,
// then call OrderSubmitter.Submit(order, signature) to send to the CLOB.
type Output struct {
	OrderID      string          // local UUID from order domain
	UnsignedHash []byte          // 32-byte EIP-712 hash for signing — NEVER signed here
	GTDExpiry    int64           // Unix timestamp of order expiration
	FeePerShare  decimal.Decimal // fee computed at current token price
}
