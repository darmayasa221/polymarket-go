package botcontainer

import (
	refreshmarkets "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/refreshmarkets"
	updateticksize "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/commands/updateticksize"
	getactivemarket "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/getactivemarket"
	ismarkettradeable "github.com/darmayasa221/polymarket-go/internal/applications/marketwatch/queries/ismarkettradeable"
	closeposition "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/closeposition"
	openposition "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/commands/openposition"
	getpnlsummary "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpnlsummary"
	getpositions "github.com/darmayasa221/polymarket-go/internal/applications/portfolio/queries/getpositions"
	recordprice "github.com/darmayasa221/polymarket-go/internal/applications/pricing/commands/recordprice"
	computefee "github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/computefee"
	getcurrentsignal "github.com/darmayasa221/polymarket-go/internal/applications/pricing/queries/getcurrentsignal"
	sharedports "github.com/darmayasa221/polymarket-go/internal/applications/shared/ports"
	cancelorder "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder"
	closewindow "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow"
	heartbeat "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat"
	placeorder "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder"
	startwindow "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow"
	tradingports "github.com/darmayasa221/polymarket-go/internal/applications/trading/ports"
	getwindowstate "github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate"
	wsmarket "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/market"
	rtds "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/rtds"
	wsuser "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/user"
)

// BotContainer holds all wired use cases and adapters for the trading bot.
// Constructed by Build; consumed by cmd/bot/.
type BotContainer struct {
	// Pricing
	RecordPrice      recordprice.UseCase
	GetCurrentSignal getcurrentsignal.UseCase
	ComputeFee       computefee.UseCase

	// Marketwatch
	RefreshMarkets    refreshmarkets.UseCase
	GetActiveMarket   getactivemarket.UseCase
	IsMarketTradeable ismarkettradeable.UseCase
	UpdateTickSize    updateticksize.UseCase

	// Trading
	PlaceOrder     placeorder.UseCase
	CancelOrder    cancelorder.UseCase
	StartWindow    startwindow.UseCase
	CloseWindow    closewindow.UseCase
	Heartbeat      heartbeat.UseCase
	GetWindowState getwindowstate.UseCase

	// Portfolio
	OpenPosition  openposition.UseCase
	ClosePosition closeposition.UseCase
	GetPositions  getpositions.UseCase
	GetPnlSummary getpnlsummary.UseCase

	// Adapters exposed to cmd/bot for direct use
	OrderSubmitter  tradingports.OrderSubmitter
	OrderRepository tradingports.OrderRepository
	FeeRateProvider sharedports.FeeRateProvider

	// WebSocket handler factories (cmd/bot calls .Start() with context)
	RTDSHandler   *rtds.Handler
	MarketHandler *wsmarket.Handler
	UserHandler   *wsuser.Handler

	// Closer releases DB and Redis connections.
	Close func() error
}
