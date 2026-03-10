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
