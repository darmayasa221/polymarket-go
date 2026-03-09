// Package json provides the JSON response presenter.
// All HTTP responses go through the presenter — never write JSON directly in handlers.
package json

import "github.com/gin-gonic/gin"

// Presenter handles all HTTP response formatting.
type Presenter struct{}

// New creates a new JSON Presenter.
func New() *Presenter { return &Presenter{} }

var _ interface {
	OK(c *gin.Context, message string, data any)
	Created(c *gin.Context, message string, data any)
	NoContent(c *gin.Context)
	Error(c *gin.Context, err error)
	Paginated(c *gin.Context, message string, data, metadata any)
} = (*Presenter)(nil)
