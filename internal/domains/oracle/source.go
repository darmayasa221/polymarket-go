// Package oracle defines price feed entities and resolution signal logic.
package oracle

// PriceSource identifies where a price reading came from.
type PriceSource string

const (
	// SourceChainlink is a Chainlink oracle price on Polygon.
	SourceChainlink PriceSource = "chainlink"
	// SourceBinance is a Binance spot price from the RTDS feed.
	SourceBinance PriceSource = "binance"
)

// IsValid returns true if the source is a supported price source.
func (s PriceSource) IsValid() bool {
	return s == SourceChainlink || s == SourceBinance
}
