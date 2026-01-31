package get

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/get"
)

// Handler maneja las peticiones de obtención de gastos
type Handler struct {
	service            get.ExpenseGetter
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de obtención de gastos
func NewHandler(service get.ExpenseGetter, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// GetExpense maneja la obtención de un gasto específico
func (h *Handler) GetExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	response, err := h.service.GetExpense(c.Request.Context(), userID, id)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordExpenseAction(
			userID,
			services.ActionViewExpenses,
			id,
			"Gasto visualizado: "+response.Description,
		)
	}

	c.JSON(http.StatusOK, response)
}
