package authentications

import (
	"strings"

	"github.com/gin-gonic/gin"

	logoutuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/logoutuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	httpconst "github.com/darmayasa221/polymarket-go/internal/interfaces/http/constants"
)

// Logout handles POST /auth/logout — invalidates the current Bearer token.
func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader(httpconst.HeaderAuthorization)
	tokenValue := strings.TrimPrefix(authHeader, httpconst.PrefixBearer)

	_, err := h.logoutUser.Execute(c.Request.Context(), logoutuserDTO.Input{
		TokenValue: tokenValue,
	})
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("logout failed", logging.FieldOperation("logout"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	h.presenter.NoContent(c)
}
