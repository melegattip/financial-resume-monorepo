package update

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	categoryUpdate "github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
)

// Handler maneja las peticiones HTTP para actualizar categorías
type Handler struct {
	service categories.CategoryService
}

// NewHandler crea una nueva instancia de Handler
func NewHandler(service categories.CategoryService) *Handler {
	return &Handler{
		service: service,
	}
}

// Handle procesa la petición para actualizar una categoría
func (h *Handler) Handle(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("User ID is required"))
		return
	}

	var request categoryUpdate.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Invalid request body: "+err.Error()))
		return
	}

	request.ID = c.Param("id")
	if request.ID == "" {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Category ID is required in URL"))
		return
	}

	// Validar que el nuevo nombre no esté vacío
	if request.NewName == "" {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("New category name is required"))
		return
	}

	response, err := h.service.Update(request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, response, "Category updated successfully")
}
