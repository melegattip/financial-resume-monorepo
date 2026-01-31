package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// IncomeHandler maneja las peticiones relacionadas con ingresos
type IncomeHandler struct {
	incomeService      incomes.IncomeService
	gamificationHelper *services.GamificationHelper
}

func NewIncomeHandler(incomeService incomes.IncomeService, gamificationHelper *services.GamificationHelper) *IncomeHandler {
	return &IncomeHandler{
		incomeService:      incomeService,
		gamificationHelper: gamificationHelper,
	}
}

// CreateIncome godoc
// @Summary Crear un nuevo ingreso
// @Description Crea un nuevo ingreso para el usuario
// @Tags incomes
// @Accept json
// @Produce json
// @Param x-caller-id header string true "ID del usuario"
// @Param income body incomes.Income true "Datos del ingreso"
// @Success 201 {object} incomes.CreateIncomeResponse
// @Failure 400 {object} errors.BadRequest
// @Failure 401 {object} errors.UnauthorizedRequest
// @Router /api/v1/incomes [post]
func (h *IncomeHandler) CreateIncome(c *gin.Context) {
	var request incomes.CreateIncomeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	request.UserID = userID

	response, err := h.incomeService.CreateIncome(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusCreated, response)
}

// GetIncome godoc
// @Summary Obtener un ingreso específico
// @Description Obtiene un ingreso por su ID
// @Tags incomes
// @Accept json
// @Produce json
// @Param x-caller-id header string true "ID del usuario"
// @Param id path string true "ID del ingreso"
// @Success 200 {object} incomes.GetIncomeResponse
// @Failure 400 {object} errors.BadRequest
// @Failure 401 {object} errors.UnauthorizedRequest
// @Failure 404 {object} errors.ResourceNotFound
// @Router /api/v1/incomes/{id} [get]
func (h *IncomeHandler) GetIncome(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	response, err := h.incomeService.GetIncome(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordIncomeAction(
			userID,
			services.ActionViewIncomes,
			id,
			"Ingreso visualizado: "+response.Description,
		)
	}

	c.JSON(http.StatusOK, response)
}

// ListIncomes godoc
// @Summary Listar ingresos
// @Description Obtiene una lista de todos los ingresos del usuario
// @Tags incomes
// @Accept json
// @Produce json
// @Param x-caller-id header string true "ID del usuario"
// @Success 200 {object} incomes.ListIncomesResponse
// @Failure 401 {object} errors.UnauthorizedRequest
// @Router /api/v1/incomes [get]
func (h *IncomeHandler) ListIncomes(c *gin.Context) {
	userID := c.GetString("user_id")
	response, err := h.incomeService.ListIncomes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, response)
}

// UpdateIncome godoc
// @Summary Actualizar un ingreso
// @Description Actualiza los datos de un ingreso existente
// @Tags incomes
// @Accept json
// @Produce json
// @Param x-caller-id header string true "ID del usuario"
// @Param id path string true "ID del ingreso"
// @Param income body incomes.Income true "Datos actualizados del ingreso"
// @Success 200 {object} incomes.UpdateIncomeResponse
// @Failure 400 {object} errors.BadRequest
// @Failure 401 {object} errors.UnauthorizedRequest
// @Failure 404 {object} errors.ResourceNotFound
// @Router /api/v1/incomes/{id} [patch]
func (h *IncomeHandler) UpdateIncome(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	var request incomes.UpdateIncomeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.incomeService.UpdateIncome(c.Request.Context(), userID, id, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ NUEVO: Registrar acción de gamificación (async, no bloquea)
	if h.gamificationHelper != nil {
		h.gamificationHelper.RecordIncomeAction(
			userID,
			services.ActionUpdateIncome,
			id,
			"Ingreso actualizado: "+response.Description,
		)
	}

	c.JSON(http.StatusOK, response)
}

// DeleteIncome godoc
// @Summary Eliminar un ingreso
// @Description Elimina un ingreso existente
// @Tags incomes
// @Accept json
// @Produce json
// @Param x-caller-id header string true "ID del usuario"
// @Param id path string true "ID del ingreso"
// @Success 204 "No Content"
// @Failure 400 {object} errors.BadRequest
// @Failure 401 {object} errors.UnauthorizedRequest
// @Failure 404 {object} errors.ResourceNotFound
// @Router /api/v1/incomes/{id} [delete]
func (h *IncomeHandler) DeleteIncome(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.incomeService.DeleteIncome(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"message": "Income deleted successfully"})
}
