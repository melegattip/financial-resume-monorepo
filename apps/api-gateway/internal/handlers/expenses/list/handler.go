package list

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/list"
)

// Handler maneja las peticiones de listado de gastos
type Handler struct {
	service            list.ExpenseLister
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de listado de gastos
func NewHandler(service list.ExpenseLister, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
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

// ListExpenses maneja el listado de gastos
func (h *Handler) ListExpenses(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	response, err := h.service.ListExpenses(c.Request.Context(), userID)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordExpenseAction(
			userID,
			services.ActionViewExpenses,
			"list",
			"Lista de gastos visualizada",
		)
	}

	c.JSON(http.StatusOK, response)
}

// ListUnpaidExpenses maneja el listado de gastos no pagados
func (h *Handler) ListUnpaidExpenses(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	response, err := h.service.ListUnpaidExpenses(c.Request.Context(), userID)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordExpenseAction(
			userID,
			services.ActionViewExpenses,
			"unpaid_list",
			"Lista de gastos pendientes visualizada",
		)
	}

	c.JSON(http.StatusOK, response)
}
