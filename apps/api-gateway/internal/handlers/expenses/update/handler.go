package update

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	expenses "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
)

// Handler maneja las peticiones de actualización de gastos
type Handler struct {
	service            expenses.ExpenseUpdater
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de actualización de gastos
func NewHandler(service expenses.ExpenseUpdater, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// UpdateExpense maneja la actualización de un gasto
func (h *Handler) UpdateExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	var request expenses.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()))
		return
	}

	response, err := h.service.UpdateExpense(c.Request.Context(), userID, id, &request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordExpenseAction(
			userID,
			services.ActionUpdateExpense,
			id,
			"Gasto actualizado: "+response.Description,
		)
	}

	c.JSON(http.StatusOK, response)
}
