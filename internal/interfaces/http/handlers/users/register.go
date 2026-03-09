package users

import (
	"github.com/gin-gonic/gin"

	adduserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/messages"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/binding"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/users/register"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// Register handles POST /users — user registration.
func (h *Handler) Register(c *gin.Context) {
	var req register.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.presenter.Error(c, binding.MapError(err))
		return
	}

	output, err := h.addUser.Execute(c.Request.Context(), adduserDTO.Input{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("register user failed", logging.FieldOperation("register"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	h.presenter.Created(c, messages.MsgUserCreated, register.Response{
		ID:        output.ID,
		Username:  output.Username,
		Email:     output.Email,
		FullName:  output.FullName,
		CreatedAt: response.JSONTime(output.CreatedAt),
	})
}
