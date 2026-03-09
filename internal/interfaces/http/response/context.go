package response

import (
	"github.com/gin-gonic/gin"
)

// ContextKeyUserID is the Gin context key used to store the authenticated user ID.
// Auth middleware sets this; handlers read it via UserIDFromContext.
const ContextKeyUserID = "user_id"

// UserIDFromContext extracts the authenticated user ID from Gin context.
func UserIDFromContext(c *gin.Context) string {
	userID, _ := c.Get(ContextKeyUserID)
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}
