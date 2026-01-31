package list

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// Handler maneja las peticiones de listado de ingresos
type Handler struct {
	service            incomes.IncomeService
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de listado de ingresos
func NewHandler(service incomes.IncomeService, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// ListIncomes maneja el listado de ingresos
func (h *Handler) ListIncomes(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("User ID is required"))
		return
	}

	incomes, err := h.service.ListIncomes(c.Request.Context(), userID)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordIncomeAction(
			userID,
			services.ActionViewIncomes,
			"list",
			"Lista de ingresos visualizada",
		)
	}

	httpUtil.SendSuccess(c, http.StatusOK, incomes, "Incomes retrieved successfully")
}
