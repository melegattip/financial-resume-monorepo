package create

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	incomesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes/create"
)

// Handler maneja las peticiones de creación de ingresos
type Handler struct {
	service            create.ServiceInterface
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia del handler de creación de ingresos
func NewHandler(service create.ServiceInterface, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// CreateIncome maneja la creación de un nuevo ingreso
func (h *Handler) CreateIncome(c *gin.Context) {
	var request incomesDomain.CreateIncomeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpUtil.BadRequest(c, errors.NewBadRequest("Invalid request body"))
		return
	}

	// Obtener userID del contexto (establecido por el middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		httpUtil.BadRequest(c, errors.NewBadRequest("User ID is required"))
		return
	}
	request.UserID = userID

	response, err := h.service.CreateIncome(c.Request.Context(), &request)
	if err != nil {
		httpUtil.HandleError(c, err)
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordIncomeAction(
			userID,
			services.ActionCreateIncome,
			response.ID,
			"Nuevo ingreso creado: "+response.Description,
		)
	}

	httpUtil.SendSuccess(c, http.StatusCreated, response, "Income created successfully")
}
