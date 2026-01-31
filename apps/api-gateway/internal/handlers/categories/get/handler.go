package get

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
)

// Handler maneja las peticiones HTTP para obtener una categoría
type Handler struct {
	service categories.CategoryService
}

// NewHandler crea una nueva instancia de Handler
func NewHandler(service categories.CategoryService) *Handler {
	return &Handler{
		service: service,
	}
}

// Handle procesa la petición para obtener una categoría
func (h *Handler) Handle(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if id == "" {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Category ID is required"))
		return
	}

	category, err := h.service.GetByID(userID, id)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, category, "Category retrieved successfully")
}
