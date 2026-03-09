package users

import (
	"github.com/gin-gonic/gin"

	listuserDTO "github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers/dto"
	"github.com/darmayasa221/polymarket-go/internal/commons/logging"
	"github.com/darmayasa221/polymarket-go/internal/commons/messages"
	"github.com/darmayasa221/polymarket-go/internal/domains/shared/pagination"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/handlers/users/list"
	httppagination "github.com/darmayasa221/polymarket-go/internal/interfaces/http/pagination"
	"github.com/darmayasa221/polymarket-go/internal/interfaces/http/response"
)

// List handles GET /users — list users with auto-detected pagination mode.
// If a "cursor" query param is present, cursor-based pagination is used.
// Otherwise, offset-based pagination (page + page_size) is used.
func (h *Handler) List(c *gin.Context) {
	if c.Query("cursor") != "" {
		h.listCursor(c)
	} else {
		h.listOffset(c)
	}
}

func (h *Handler) listOffset(c *gin.Context) {
	input := listuserDTO.Input{
		Mode:         pagination.ModeOffset,
		OffsetParams: httppagination.ParseOffset(c),
	}

	output, err := h.listUsers.ExecuteOffset(c.Request.Context(), input)
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("list users failed", logging.FieldOperation("list"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	items := mapUserItems(output.Users)
	meta := httppagination.NewOffsetMeta(output.Page, output.PageSize, output.TotalItems, output.TotalPages)
	h.presenter.Paginated(c, messages.MsgUsersListed, items, meta)
}

func (h *Handler) listCursor(c *gin.Context) {
	input := listuserDTO.Input{
		Mode:         pagination.ModeCursor,
		CursorParams: httppagination.ParseCursor(c),
	}

	output, err := h.listUsers.ExecuteCursor(c.Request.Context(), input)
	if err != nil {
		logging.FromContext(c.Request.Context(), h.logger).Error("list users failed", logging.FieldOperation("list"), logging.FieldLayer("handler"), logging.FieldError(err))
		h.presenter.Error(c, err)
		return
	}

	items := mapUserItems(output.Users)
	meta := httppagination.NewCursorMeta(output.NextCursor, output.PrevCursor, output.HasNext, output.HasPrev)
	h.presenter.Paginated(c, messages.MsgUsersListed, items, meta)
}

func mapUserItems(users []listuserDTO.UserItem) []list.UserItem {
	items := make([]list.UserItem, len(users))
	for i, u := range users {
		items[i] = list.UserItem{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			CreatedAt: response.JSONTime(u.CreatedAt),
		}
	}
	return items
}
