// Package pagination provides HTTP-specific pagination parameter parsing.
// Domain pagination types live in domains/shared/pagination — this only parses HTTP params.
package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
)

// ParseOffset extracts offset pagination params from an HTTP request.
func ParseOffset(c *gin.Context) pagination.OffsetParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", defaultPageStr))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", defaultPageSizeStr))
	return pagination.NewOffsetParams(page, pageSize)
}

// ParseCursor extracts cursor pagination params from an HTTP request.
func ParseCursor(c *gin.Context) pagination.CursorParams {
	cursor := c.Query("cursor")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", defaultPageSizeStr))
	forward := c.DefaultQuery("direction", defaultDirectionStr) == defaultDirectionStr
	return pagination.NewCursorParams(cursor, pageSize, forward)
}
