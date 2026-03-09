package authentications

import (
	"github.com/gin-gonic/gin"

	loginuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/authentications/commands/loginuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/messages"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/authentications/login"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/binding"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// Login handles POST /auth/login — authenticates a user and returns a token pair.
func (h *Handler) Login(c *gin.Context) {
	var req login.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.presenter.Error(c, binding.MapError(err))
		return
	}

	output, err := h.loginUser.Execute(c.Request.Context(), loginuserDTO.Input{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("login failed", logging.FieldOperation("login"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	h.presenter.OK(c, messages.MsgLoginSuccess, login.Response{
		AccessToken:           output.AccessToken,
		RefreshToken:          output.RefreshToken,
		AccessTokenExpiresAt:  response.JSONTime(output.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: response.JSONTime(output.RefreshTokenExpiresAt),
		ID:                    output.UserID,
	})
}
