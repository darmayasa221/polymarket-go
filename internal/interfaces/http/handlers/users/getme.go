package users

import (
	"github.com/gin-gonic/gin"

	getuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/messages"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/users/getme"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// GetMe handles GET /users/me — retrieve authenticated user's own profile.
func (h *Handler) GetMe(c *gin.Context) {
	userID := response.UserIDFromContext(c)
	if userID == "" {
		h.presenter.Error(c, unauthorizedError())
		return
	}

	output, err := h.getUser.Execute(c.Request.Context(), getuserDTO.Input{UserID: userID})
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("get me failed", logging.FieldOperation("getme"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	h.presenter.OK(c, messages.MsgUserFound, getme.Response{
		ID:        output.ID,
		Username:  output.Username,
		Email:     output.Email,
		FullName:  output.FullName,
		CreatedAt: response.JSONTime(output.CreatedAt),
		UpdatedAt: response.JSONTime(output.UpdatedAt),
	})
}
