package delete

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/delete"
)

// Handler maneja las peticiones de eliminación de gastos
type Handler struct {
	service            delete.ExpenseDeleter
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de eliminación de gastos
func NewHandler(service delete.ExpenseDeleter, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// DeleteExpense maneja la eliminación de un gasto
func (h *Handler) DeleteExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.service.DeleteExpense(c.Request.Context(), userID, id); err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordExpenseAction(
			userID,
			services.ActionDeleteExpense,
			id,
			"Gasto eliminado",
		)
	}

	c.JSON(http.StatusNoContent, nil)
}
