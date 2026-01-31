package create

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
)

// Handler maneja las peticiones HTTP para crear categorías
type Handler struct {
	service categories.CategoryService
}

// NewHandler crea una nueva instancia de Handler
func NewHandler(service categories.CategoryService) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateCategoryRequestBody representa solo los datos necesarios en el body
type CreateCategoryRequestBody struct {
	Name string `json:"name" binding:"required"`
}

// getUserIDFromContext extrae el user_id del contexto JWT
func (h *Handler) getUserIDFromContext(c *gin.Context) (string, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", coreErrors.NewUnauthorizedRequest("Usuario no autenticado")
	}

	// Convertir a string según el tipo
	switch v := userIDInterface.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case string:
		return v, nil
	default:
		return "", coreErrors.NewBadRequest("Formato de user_id inválido")
	}
}

// Handle procesa la petición para crear una categoría
func (h *Handler) Handle(c *gin.Context) {
	var requestBody CreateCategoryRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Invalid request body"))
		return
	}

	// Obtener userID desde el contexto JWT (consistente con otros handlers)
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	if requestBody.Name == "" {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Name is required"))
		return
	}

	// Crear el request completo con el userID del contexto JWT
	request := categories.CreateCategoryRequest{
		UserID: userID,
		Name:   requestBody.Name,
	}

	category, err := h.service.Create(request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusCreated, category, "Category created successfully")
}
