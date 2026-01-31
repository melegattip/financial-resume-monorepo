package list

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
)

// Handler maneja las peticiones HTTP para listar categorías
type Handler struct {
	service categories.CategoryService
}

// NewHandler crea una nueva instancia de Handler
func NewHandler(service categories.CategoryService) *Handler {
	return &Handler{
		service: service,
	}
}

// getUserIDFromContext extrae el user_id del contexto JWT
func (h *Handler) getUserIDFromContext(c *gin.Context) (string, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", errors.NewUnauthorizedRequest("Usuario no autenticado")
	}

	// Convertir a string según el tipo
	switch v := userIDInterface.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case string:
		return v, nil
	default:
		return "", errors.NewBadRequest("Formato de user_id inválido")
	}
}

// Handle procesa la petición para listar categorías
func (h *Handler) Handle(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	categories, err := h.service.List(userID)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	httpUtil.SendSuccess(c, http.StatusOK, categories, "Categories retrieved successfully")
}
