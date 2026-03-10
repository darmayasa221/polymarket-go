package dto

import "github.com/shopspring/decimal"

// Input holds the parameters for placing a limit order on the CLOB.
type Input struct {
	Asset         string          // "btc" | "eth" | "sol" | "xrp"
	Outcome       string          // "Up" | "Down"
	Side          string          // "buy" | "sell"
	Price         decimal.Decimal // token price 0.01–0.99
	Size          decimal.Decimal // shares to trade (>= 5, minimum order size)
	TokenID       string          // ERC1155 token ID from GetActiveMarket
	FeeRateBps    uint64          // live fee rate — NEVER hardcode; fetched from CLOB /fee-rate
	FunderAddress string          // EOA wallet address (from env)
}
