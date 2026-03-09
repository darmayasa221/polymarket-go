// Package server provides the HTTP server implementation.
package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/config"
	"github.com/darmayasa221/polymarket-go/internal/infrastructures/container"
	authhandler "github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/authentications"
	usershandler "github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/users"
	bodysizemw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/bodysize"
	corsmw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/cors"
	loggermw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/logger"
	recoverymw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/recovery"
	securitymw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/security"
	timeoutmw "github.com/darmayasa221/polymarket-go/internal/interfaces/http/middlewares/timeout"
	jsonpresenter "github.com/darmayasa221/polymarket-go/internal/interfaces/http/presenters/json"
	authroutes "github.com/darmayasa221/polymarket-go/internal/interfaces/http/routes/authentications"
	healthroutes "github.com/darmayasa221/polymarket-go/internal/interfaces/http/routes/health"
	userroutes "github.com/darmayasa221/polymarket-go/internal/interfaces/http/routes/users"
)

// Server wraps the HTTP server and gin engine.
type Server struct {
	engine *gin.Engine
	server *http.Server
	logger *logging.Logger
}

// New builds the gin engine with the full middleware chain and registers all routes.
func New(cfg *config.Config, c *container.Container, logger *logging.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Middleware chain (Chain of Responsibility).
	// Order: recovery → cors → logger → security → bodysize → timeout
	engine.Use(
		recoverymw.New(logger),
		corsmw.New(cfg.HTTP.AllowedOrigins),
		loggermw.New(logger),
		securitymw.New(),
		bodysizemw.New(cfg.HTTP.MaxBodyBytes),
		timeoutmw.New(cfg.HTTP.RequestTimeout),
	)

	presenter := jsonpresenter.New()

	// Health routes on root (no /api/v1 prefix).
	healthroutes.Register(engine.Group(""), c.HealthChecker)

	// API v1 group.
	v1 := engine.Group("/api/v1")

	usersH := usershandler.New(c.AddUser, c.GetUser, c.ListUsers, presenter, logger)
	authH := authhandler.New(c.LoginUser, c.LogoutUser, c.RefreshAuth, presenter, logger)

	userroutes.Register(v1, usersH, c.TokenManager)
	authroutes.Register(v1, authH)

	return &Server{
		engine: engine,
		server: &http.Server{
			Addr:         ":" + cfg.HTTP.Port,
			Handler:      engine,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
		logger: logger,
	}
}

// Start begins listening and serving HTTP requests.
// Blocks until the server is shut down or encounters an error.
// Returns http.ErrServerClosed on graceful shutdown.
func (s *Server) Start() error {
	s.logger.Info("http server starting", logging.FieldLayer("server"))
	return s.server.ListenAndServe()
}

// Shutdown gracefully drains active connections within the given context deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
