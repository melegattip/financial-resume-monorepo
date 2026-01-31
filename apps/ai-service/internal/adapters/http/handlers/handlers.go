package handlers

import (
	"net/http"
	"strconv"

	"github.com/financial-ai-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// Handlers contiene todos los handlers del servicio
type Handlers struct {
	analysisUseCase ports.AIAnalysisPort
	purchaseUseCase ports.PurchaseDecisionPort
	creditUseCase   ports.CreditAnalysisPort
}

// NewHandlers crea una nueva instancia de handlers
func NewHandlers(
	analysisUseCase ports.AIAnalysisPort,
	purchaseUseCase ports.PurchaseDecisionPort,
	creditUseCase ports.CreditAnalysisPort,
) *Handlers {
	return &Handlers{
		analysisUseCase: analysisUseCase,
		purchaseUseCase: purchaseUseCase,
		creditUseCase:   creditUseCase,
	}
}

// SetupRoutes configura todas las rutas del servicio
func SetupRoutes(
	router *gin.Engine,
	analysisUseCase ports.AIAnalysisPort,
	purchaseUseCase ports.PurchaseDecisionPort,
	creditUseCase ports.CreditAnalysisPort,
) {
	handlers := NewHandlers(analysisUseCase, purchaseUseCase, creditUseCase)

	// Grupo de rutas para AI
	ai := router.Group("/api/v1/ai")
	{
		// Análisis financiero
		ai.POST("/health-analysis", handlers.AnalyzeFinancialHealth)
		ai.POST("/insights", handlers.GenerateInsights)

		// Decisiones de compra
		ai.POST("/can-i-buy", handlers.CanIBuy)
		ai.POST("/alternatives", handlers.SuggestAlternatives)

		// Análisis crediticio
		ai.POST("/credit-plan", handlers.GenerateCreditPlan)
		ai.POST("/credit-score", handlers.CalculateCreditScore)
	}

	// Health check
	router.GET("/health", handlers.HealthCheck)
}

// AnalyzeFinancialHealth maneja el análisis de salud financiera
func (h *Handlers) AnalyzeFinancialHealth(c *gin.Context) {
	var request ports.FinancialAnalysisData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	analysis, err := h.analysisUseCase.AnalyzeFinancialHealth(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error analyzing financial health",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analysis,
		"source":  "ai_microservice",
	})
}

// GenerateInsights maneja la generación de insights financieros
func (h *Handlers) GenerateInsights(c *gin.Context) {
	var request ports.FinancialAnalysisData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	insights, err := h.analysisUseCase.GenerateInsights(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error generating insights",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    insights,
		"source":  "ai_microservice",
	})
}

// CanIBuy maneja el análisis de decisiones de compra
func (h *Handlers) CanIBuy(c *gin.Context) {
	var request ports.PurchaseAnalysisRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	decision, err := h.purchaseUseCase.CanIBuy(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error analyzing purchase decision",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    decision,
		"source":  "ai_microservice",
	})
}

// SuggestAlternatives maneja la sugerencia de alternativas de compra
func (h *Handlers) SuggestAlternatives(c *gin.Context) {
	var request ports.PurchaseAnalysisRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	alternatives, err := h.purchaseUseCase.SuggestAlternatives(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error generating alternatives",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    alternatives,
		"source":  "ai_microservice",
	})
}

// GenerateCreditPlan maneja la generación de planes de mejora crediticia
func (h *Handlers) GenerateCreditPlan(c *gin.Context) {
	var request ports.FinancialAnalysisData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	plan, err := h.creditUseCase.GenerateImprovementPlan(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error generating credit improvement plan",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
		"source":  "ai_microservice",
	})
}

// CalculateCreditScore maneja el cálculo de score crediticio
func (h *Handlers) CalculateCreditScore(c *gin.Context) {
	var request ports.FinancialAnalysisData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	score, err := h.creditUseCase.CalculateCreditScore(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error calculating credit score",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"score":         score,
			"user_id":       request.UserID,
			"calculated_at": "now",
		},
		"source": "ai_microservice",
	})
}

// HealthCheck maneja el endpoint de health check
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "Financial AI Service",
		"version": "1.0.0",
		"port":    c.Request.Header.Get("X-Forwarded-Port"),
	})
}

// Utility functions

// getFloatParam obtiene un parámetro float de la query string
func getFloatParam(c *gin.Context, key string, defaultValue float64) float64 {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getIntParam obtiene un parámetro int de la query string
func getIntParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getBoolParam obtiene un parámetro bool de la query string
func getBoolParam(c *gin.Context, key string, defaultValue bool) bool {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
