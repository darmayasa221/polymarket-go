package main

import (
	"context"

	appconst "github.com/darmayasa221/polymarket-go/internal/commons/constants/app"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container"
)

type app struct {
	servers []Server
	logger  *logging.Logger
}

func build() (*app, error) {
	cfg, err := config.Load(".env")
	if err != nil {
		return nil, err
	}

	logger, err := logging.New(cfg.App.LogLevel)
	if err != nil {
		return nil, err
	}

	c, err := container.Build(cfg, logger)
	if err != nil {
		return nil, err
	}

	servers := buildServers(cfg, c, logger)
	return &app{servers: servers, logger: logger}, nil
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), appconst.DefaultShutdownTimeout)
	defer cancel()

	for _, srv := range a.servers {
		if err := srv.Shutdown(ctx); err != nil {
			a.logger.Error("shutdown error", logging.FieldError(err))
		}
	}
	_ = a.logger.Sync()
}
