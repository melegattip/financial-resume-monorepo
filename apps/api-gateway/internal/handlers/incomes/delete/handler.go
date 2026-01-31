package delete

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes/delete"
)

// Handler maneja las peticiones HTTP para eliminar ingresos
type Handler struct {
	service            delete.IncomeDeleter
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler
func NewHandler(service delete.IncomeDeleter, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// Handle procesa la petición HTTP para eliminar un ingreso
func (h *Handler) Handle(c *gin.Context) {
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

	err := h.service.DeleteIncome(c.Request.Context(), userID, id)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordIncomeAction(
			userID,
			services.ActionDeleteIncome,
			id,
			"Ingreso eliminado",
		)
	}

	httpUtil.SendSuccess(c, http.StatusOK, nil, "Income deleted successfully")
}
