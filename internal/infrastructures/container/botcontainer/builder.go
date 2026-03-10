package botcontainer

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

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
	cancelorder "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/cancelorder"
	closewindow "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/closewindow"
	heartbeat "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/heartbeat"
	placeorder "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/placeorder"
	startwindow "github.com/darmayasa221/polymarket-go/internal/applications/trading/commands/startwindow"
	getwindowstate "github.com/darmayasa221/polymarket-go/internal/applications/trading/queries/getwindowstate"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
	gammaclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/gamma"
	marketpgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/marketwatch/postgres"
	portfoliopgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/portfolio/postgres"
	pricingpgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/pricing/postgres"
	tradingpgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/trading/postgres"
	tradingredisrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/trading/redis"
	wsmarket "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/market"
	rtds "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/rtds"
	wsuser "github.com/darmayasa221/polymarket-go/internal/infrastructures/ws/user"
)

// Config holds the subset of app config needed by the bot container.
type Config struct {
	PostgresDSN   string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	CLOBBaseURL   string
	CLOBAddress   string
	CLOBAPIKey    string
	CLOBAPISecret string
	CLOBAPIPass   string
	GammaBaseURL  string
}

// FromAppConfig converts the full app config to a bot Config.
func FromAppConfig(cfg *config.Config) Config {
	return Config{
		PostgresDSN:   cfg.PostgreSQL.DSN,
		RedisAddress:  cfg.Cache.Address,
		RedisPassword: cfg.Cache.Password,
		RedisDB:       cfg.Cache.DB,
		CLOBBaseURL:   cfg.CLOB.BaseURL,
		CLOBAddress:   cfg.CLOB.Address,
		CLOBAPIKey:    cfg.CLOB.APIKey,
		CLOBAPISecret: cfg.CLOB.APISecret,
		CLOBAPIPass:   cfg.CLOB.APIPassphrase,
		GammaBaseURL:  cfg.Gamma.BaseURL,
	}
}

// Build wires all Polymarket adapters and use cases, returning a BotContainer.
// On error, any opened connections are closed before returning.
func Build(cfg Config) (*BotContainer, error) {
	// --- infrastructure ---
	db, err := pgdb.New(pgdb.Config{DSN: cfg.PostgresDSN})
	if err != nil {
		return nil, fmt.Errorf("botcontainer: postgres: %w", err)
	}

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	clobClient := clob.NewClient(clob.Config{
		BaseURL:       cfg.CLOBBaseURL,
		Address:       cfg.CLOBAddress,
		APIKey:        cfg.CLOBAPIKey,
		APISecret:     cfg.CLOBAPISecret,
		APIPassphrase: cfg.CLOBAPIPass,
	})

	gammaSource := gammaclient.NewMarketSource(gammaclient.Config{BaseURL: cfg.GammaBaseURL})

	// --- repositories ---
	priceRepo := pricingpgrepo.New(db)
	marketRepo := marketpgrepo.New(db)
	orderRepo := tradingpgrepo.New(db)
	positionRepo := portfoliopgrepo.New(db)
	store := tradingredisrepo.New(redisClient)

	// --- CLOB adapters ---
	feeProvider := clob.NewFeeRateProvider(clobClient)
	submitter := clob.NewOrderSubmitter(clobClient)
	sender := clob.NewHeartbeatSender(clobClient)
	balanceProvider := clob.NewBalanceProvider(clobClient)

	return &BotContainer{
		RecordPrice:      recordprice.New(priceRepo),
		GetCurrentSignal: getcurrentsignal.New(priceRepo),
		ComputeFee:       computefee.New(),

		RefreshMarkets:    refreshmarkets.New(gammaSource, marketRepo),
		GetActiveMarket:   getactivemarket.New(marketRepo),
		IsMarketTradeable: ismarkettradeable.New(marketRepo),
		UpdateTickSize:    updateticksize.New(marketRepo),

		PlaceOrder:     placeorder.New(orderRepo, store),
		CancelOrder:    cancelorder.New(orderRepo, submitter),
		StartWindow:    startwindow.New(store),
		CloseWindow:    closewindow.New(store, orderRepo),
		Heartbeat:      heartbeat.New(sender),
		GetWindowState: getwindowstate.New(store),

		OpenPosition:  openposition.New(positionRepo),
		ClosePosition: closeposition.New(positionRepo),
		GetPositions:  getpositions.New(positionRepo),
		GetPnlSummary: getpnlsummary.New(positionRepo),

		OrderSubmitter:  submitter,
		OrderRepository: orderRepo,
		FeeRateProvider: feeProvider,
		BalanceProvider: balanceProvider,

		RTDSHandler:   rtds.New(rtds.RTDSEndpoint),
		MarketHandler: wsmarket.New(wsmarket.MarketEndpoint),
		UserHandler:   wsuser.New(wsuser.UserEndpoint),

		Close: func() error {
			_ = redisClient.Close()
			return db.Close()
		},

		pgDB: db,
	}, nil
}

// RunMigration runs the PostgreSQL schema migration.
// Uses CREATE TABLE IF NOT EXISTS — safe to call on every startup.
func (bc *BotContainer) RunMigration(_ context.Context) error {
	return pgdb.RunMigrations(bc.pgDB.DB())
}
