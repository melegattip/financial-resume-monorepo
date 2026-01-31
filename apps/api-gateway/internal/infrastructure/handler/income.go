package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// IncomeHandler maneja las peticiones HTTP relacionadas con ingresos
type IncomeHandler struct {
	service incomes.IncomeService
}

// NewIncomeHandler crea una nueva instancia del handler de ingresos
func NewIncomeHandler(service incomes.IncomeService) *IncomeHandler {
	return &IncomeHandler{
		service: service,
	}
}

// CreateIncome maneja la creación de un nuevo ingreso
func (h *IncomeHandler) CreateIncome(c *gin.Context) {
	var request incomes.CreateIncomeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar campos requeridos
	if request.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description is required"})
		return
	}
	if request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	response, err := h.service.CreateIncome(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetIncome maneja la obtención de un ingreso por su ID
func (h *IncomeHandler) GetIncome(c *gin.Context) {
	userID := c.Param("user_id")
	incomeID := c.Param("id")

	response, err := h.service.GetIncome(c.Request.Context(), userID, incomeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListIncomes maneja la obtención de todos los ingresos de un usuario
func (h *IncomeHandler) ListIncomes(c *gin.Context) {
	userID := c.Param("user_id")

	response, err := h.service.ListIncomes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateIncome maneja la actualización de un ingreso
func (h *IncomeHandler) UpdateIncome(c *gin.Context) {
	userID := c.Param("user_id")
	incomeID := c.Param("id")

	var request incomes.UpdateIncomeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateIncome(c.Request.Context(), userID, incomeID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteIncome maneja la eliminación de un ingreso
func (h *IncomeHandler) DeleteIncome(c *gin.Context) {
	userID := c.Param("user_id")
	incomeID := c.Param("id")

	err := h.service.DeleteIncome(c.Request.Context(), userID, incomeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
