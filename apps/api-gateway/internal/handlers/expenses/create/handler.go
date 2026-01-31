package create

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/create"
)

// Handler maneja las peticiones de creación de gastos
type Handler struct {
	service            create.ServiceInterface
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de creación de gastos
func NewHandler(service create.ServiceInterface, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// CreateExpense maneja la creación de un nuevo gasto
func (h *Handler) CreateExpense(c *gin.Context) {
	var request expensesDomain.CreateExpenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Invalid request body"))
		return
	}

	// Obtener userID del contexto (establecido por el middleware)
	userID := c.GetString("user_id")
	request.UserID = userID

	response, err := h.service.CreateExpense(c.Request.Context(), &request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordExpenseAction(
			userID,
			services.ActionCreateExpense,
			response.ID,
			"Nuevo gasto creado: "+response.Description,
		)
	}

	httpUtil.SendSuccess(c, http.StatusCreated, response, "Expense created successfully")
}
