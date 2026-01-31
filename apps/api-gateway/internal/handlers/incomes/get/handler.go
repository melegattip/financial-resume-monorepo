package get

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// Handler maneja las peticiones de obtención de un ingreso
type Handler struct {
	service            incomes.IncomeService
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de obtención de un ingreso
func NewHandler(service incomes.IncomeService, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// GetIncome maneja la obtención de un ingreso
func (h *Handler) GetIncome(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("User ID is required"))
		return
	}

	id := c.Param("id")
	if id == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("Income ID is required"))
		return
	}

	income, err := h.service.GetIncome(c.Request.Context(), userID, id)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordIncomeAction(
			userID,
			services.ActionViewIncomes,
			id,
			"Ingreso visualizado: "+income.Description,
		)
	}

	httpUtil.SendSuccess(c, http.StatusOK, income, "Income retrieved successfully")
}
