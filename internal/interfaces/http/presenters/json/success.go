package json

import (
	"github.com/gin-gonic/gin"

	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// successBody is the standard success response envelope.
type successBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK sends a 200 success response.
func (p *Presenter) OK(c *gin.Context, message string, data any) {
	c.JSON(httpconst.StatusOK, successBody{Success: true, Message: message, Data: data})
}

// Created sends a 201 created response.
func (p *Presenter) Created(c *gin.Context, message string, data any) {
	c.JSON(httpconst.StatusCreated, successBody{Success: true, Message: message, Data: data})
}

// NoContent sends a 204 no content response.
func (p *Presenter) NoContent(c *gin.Context) {
	c.Status(httpconst.StatusNoContent)
}
