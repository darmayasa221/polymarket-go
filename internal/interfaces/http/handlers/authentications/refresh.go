package authentications

import (
	"github.com/gin-gonic/gin"

	refreshauthDTO "github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/refreshauth/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/messages"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/authentications/refresh"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/binding"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// Refresh handles POST /auth/refresh — exchanges a refresh token for a new token pair.
func (h *Handler) Refresh(c *gin.Context) {
	var req refresh.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.presenter.Error(c, binding.MapError(err))
		return
	}

	output, err := h.refreshAuth.Execute(c.Request.Context(), refreshauthDTO.Input{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("token refresh failed", logging.FieldOperation("refresh"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	h.presenter.OK(c, messages.MsgRefreshSuccess, refresh.Response{
		AccessToken:           output.AccessToken,
		RefreshToken:          output.RefreshToken,
		AccessTokenExpiresAt:  response.JSONTime(output.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: response.JSONTime(output.RefreshTokenExpiresAt),
	})
}
