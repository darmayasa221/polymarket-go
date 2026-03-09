// Package authentications provides the HTTP handler for the authentications domain.
package authentications

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	jsonpresenter "github.com/darmayasa221/polymarket-go/internal/interfaces/http/presenters/json"
)

// Handler handles HTTP requests for the authentications domain.
// Each action is in its own file — this file contains only the struct and constructor.
type Handler struct {
	loginUser   loginuser.UseCase
	logoutUser  logoutuser.UseCase
	refreshAuth refreshauth.UseCase
	presenter   *jsonpresenter.Presenter
	logger      *logging.Logger
}

// New creates a new authentications Handler.
func New(
	loginUser loginuser.UseCase,
	logoutUser logoutuser.UseCase,
	refreshAuth refreshauth.UseCase,
	presenter *jsonpresenter.Presenter,
	logger *logging.Logger,
) *Handler {
	return &Handler{
		loginUser:   loginUser,
		logoutUser:  logoutUser,
		refreshAuth: refreshAuth,
		presenter:   presenter,
		logger:      logger,
	}
}
