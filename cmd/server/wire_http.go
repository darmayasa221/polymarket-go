package main

import (
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container"
	httpserver "github.com/darmayasa221/polymarket-go/internal/interfaces/http/server"
)

func buildServers(cfg *config.Config, c *container.Container, logger *logging.Logger) []Server {
	return []Server{
		httpserver.New(cfg, c, logger),
	}
}
