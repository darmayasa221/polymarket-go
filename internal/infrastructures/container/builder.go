package container

import (
	"fmt"

	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container/providers"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container/providers/usecases"
)

// Build wires every provider in dependency order and returns a fully populated Container.
// On any error, already-opened resources are closed before returning.
// The logger parameter is reserved for provider startup diagnostics (e.g. "database connected").
func Build(cfg *config.Config, _ *logging.Logger) (*Container, error) {
	db, err := providers.ProvideDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("container: build: %w", err)
	}

	cache, err := providers.ProvideCache(cfg)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("container: build: %w", err)
	}

	enc := providers.ProvideEncryption(cfg)
	tm := providers.ProvideTokenManager(cfg)

	userRepo := providers.ProvideUserRepository(db, cache)
	authRepo := providers.ProvideAuthRepository(cache)

	addUser := usecases.ProvideAddUser(userRepo, enc)
	getUser := usecases.ProvideGetUser(userRepo)
	listUsers := usecases.ProvideListUsers(userRepo)
	loginUser := usecases.ProvideLoginUser(userRepo, authRepo, enc, tm)
	logoutUser := usecases.ProvideLogoutUser(authRepo)
	refreshAuth := usecases.ProvideRefreshAuth(authRepo, tm)

	healthChecker := providers.ProvideHealthChecker(db, cache)

	return &Container{
		AddUser:       addUser,
		GetUser:       getUser,
		ListUsers:     listUsers,
		LoginUser:     loginUser,
		LogoutUser:    logoutUser,
		RefreshAuth:   refreshAuth,
		TokenManager:  tm,
		HealthChecker: healthChecker,
	}, nil
}
