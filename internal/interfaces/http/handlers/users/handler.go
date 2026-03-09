// Package users provides the HTTP handler for the users domain.
package users

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	jsonpresenter "github.com/darmayasa221/polymarket-go/internal/interfaces/http/presenters/json"
)

// Handler handles HTTP requests for the users domain.
// Each action is in its own file — this file contains only the struct and constructor.
type Handler struct {
	addUser   adduser.UseCase
	getUser   getuser.UseCase
	listUsers listusers.UseCase
	presenter *jsonpresenter.Presenter
	logger    *logging.Logger
}

// New creates a new users Handler.
func New(
	addUser adduser.UseCase,
	getUser getuser.UseCase,
	listUsers listusers.UseCase,
	presenter *jsonpresenter.Presenter,
	logger *logging.Logger,
) *Handler {
	return &Handler{
		addUser:   addUser,
		getUser:   getUser,
		listUsers: listUsers,
		presenter: presenter,
		logger:    logger,
	}
}
