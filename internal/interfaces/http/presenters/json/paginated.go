package json

import (
	"github.com/gin-gonic/gin"

	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// paginatedBody is the standard paginated response envelope.
type paginatedBody struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
	Metadata any    `json:"metadata"`
}

// Paginated sends a paginated success response with metadata.
func (p *Presenter) Paginated(c *gin.Context, message string, data, metadata any) {
	c.JSON(httpconst.StatusOK, paginatedBody{
		Success:  true,
		Message:  message,
		Data:     data,
		Metadata: metadata,
	})
}
