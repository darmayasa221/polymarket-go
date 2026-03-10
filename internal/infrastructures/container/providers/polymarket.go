// Package providers wires polymarket infrastructure adapters into the container.
package providers

import (
	goredis "github.com/redis/go-redis/v9"

	"github.com/darmayasa221/polymarket-go/internal/infrastructures/clob"
	pgdb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/postgres"
	gammaclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/gamma"
	marketpgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/marketwatch/postgres"
	portfoliopgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/portfolio/postgres"
	pricingpgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/pricing/postgres"
	tradingpgrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/trading/postgres"
	tradingredisrepo "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/trading/redis"
)

// ProvidePostgresDB opens a PostgreSQL connection from the given DSN.
func ProvidePostgresDB(dsn string) (*pgdb.DB, error) {
	return pgdb.New(pgdb.Config{DSN: dsn})
}

// ProvideFeeRateProvider wires the CLOB fee-rate provider.
func ProvideFeeRateProvider(client *clob.Client) *clob.FeeRateProvider {
	return clob.NewFeeRateProvider(client)
}

// ProvideOrderSubmitter wires the CLOB order submitter.
func ProvideOrderSubmitter(client *clob.Client) *clob.OrderSubmitter {
	return clob.NewOrderSubmitter(client)
}

// ProvideHeartbeatSender wires the CLOB heartbeat sender.
func ProvideHeartbeatSender(client *clob.Client) *clob.HeartbeatSender {
	return clob.NewHeartbeatSender(client)
}

// ProvidePricingRepository wires the PostgreSQL price repository.
func ProvidePricingRepository(db *pgdb.DB) *pricingpgrepo.Repository {
	return pricingpgrepo.New(db)
}

// ProvideMarketRepository wires the PostgreSQL market repository.
func ProvideMarketRepository(db *pgdb.DB) *marketpgrepo.Repository {
	return marketpgrepo.New(db)
}

// ProvideOrderRepository wires the PostgreSQL order repository.
func ProvideOrderRepository(db *pgdb.DB) *tradingpgrepo.Repository {
	return tradingpgrepo.New(db)
}

// ProvidePositionRepository wires the PostgreSQL position repository.
func ProvidePositionRepository(db *pgdb.DB) *portfoliopgrepo.Repository {
	return portfoliopgrepo.New(db)
}

// ProvideWindowStateStore wires the Redis WindowStateStore.
func ProvideWindowStateStore(client *goredis.Client) *tradingredisrepo.Store {
	return tradingredisrepo.New(client)
}

// ProvideCLOBClient wires the CLOB HTTP client from config.
func ProvideCLOBClient(baseURL, address, apiKey, apiSecret, apiPassphrase string) *clob.Client {
	return clob.NewClient(clob.Config{
		BaseURL:       baseURL,
		Address:       address,
		APIKey:        apiKey,
		APISecret:     apiSecret,
		APIPassphrase: apiPassphrase,
	})
}

// ProvideGammaMarketSource wires the Gamma API MarketSource from config.
func ProvideGammaMarketSource(baseURL string) *gammaclient.MarketSource {
	return gammaclient.NewMarketSource(gammaclient.Config{BaseURL: baseURL})
}
