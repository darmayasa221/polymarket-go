package users

import (
	"github.com/gin-gonic/gin"

	getuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/messages"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/binding"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/users/getbyid"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// GetByID handles GET /users/:id — retrieve user by ID.
func (h *Handler) GetByID(c *gin.Context) {
	var params getbyid.URIParams
	if err := c.ShouldBindUri(&params); err != nil {
		h.presenter.Error(c, binding.MapError(err))
		return
	}

	output, err := h.getUser.Execute(c.Request.Context(), getuserDTO.Input{UserID: params.ID})
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("get user by id failed", logging.FieldOperation("getbyid"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	h.presenter.OK(c, messages.MsgUserFound, getbyid.Response{
		ID:        output.ID,
		Username:  output.Username,
		Email:     output.Email,
		FullName:  output.FullName,
		CreatedAt: response.JSONTime(output.CreatedAt),
		UpdatedAt: response.JSONTime(output.UpdatedAt),
	})
}
